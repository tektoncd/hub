// Copyright 2021 Google LLC
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     https://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package spannerdriver

import (
	"bytes"
	"context"
	"crypto/sha256"
	"encoding/gob"
	"reflect"

	"cloud.google.com/go/spanner"
	sppb "cloud.google.com/go/spanner/apiv1/spannerpb"
	"google.golang.org/api/iterator"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/structpb"
)

var errNextAfterSTop = status.Errorf(codes.FailedPrecondition, "Next called after Stop")

// init registers the protobuf types with gob so they can be encoded.
func init() {
	gob.Register(structpb.Value_BoolValue{})
	gob.Register(structpb.Value_NumberValue{})
	gob.Register(structpb.Value_StringValue{})
	gob.Register(structpb.Value_NullValue{})
	gob.Register(structpb.Value_ListValue{})
	gob.Register(structpb.Value_StructValue{})
}

// checksumRowIterator implements rowIterator and keeps track of a running
// checksum for all results that have been seen during the iteration of the
// results. This checksum can be used to verify whether a retry returned the
// same results as the initial attempt or not.
type checksumRowIterator struct {
	*spanner.RowIterator
	metadata *sppb.ResultSetMetadata

	ctx  context.Context
	tx   *readWriteTransaction
	stmt spanner.Statement
	// nc (nextCount) indicates the number of times that next has been called
	// on the iterator. Next() will be called the same number of times during
	// a retry.
	nc int64
	// stopped indicates whether the original iterator was stopped. If it was,
	// the iterator that is created during a retry should also be stopped after
	// the retry has finished.
	stopped bool

	// checksum contains the current checksum for the results that have been
	// seen. It is calculated as a SHA256 checksum over all rows that so far
	// have been returned.
	checksum *[32]byte
	buffer   *bytes.Buffer
	enc      *gob.Encoder

	// errIndex and err indicate any error and the index in the result set
	// where the error occurred.
	errIndex int64
	err      error
}

func (it *checksumRowIterator) Next() (row *spanner.Row, err error) {
	if it.stopped {
		return nil, errNextAfterSTop
	}
	err = it.tx.runWithRetry(it.ctx, func(ctx context.Context) error {
		row, err = it.RowIterator.Next()
		// spanner.ErrCode returns codes.Ok for nil errors.
		if spanner.ErrCode(err) == codes.Aborted {
			return err
		}
		if err != nil {
			// Register the error that we received and the row where we
			// received it. This will in almost all cases be the first row
			// when the query fails, or the last row when the iterator
			// returns iterator.Done. It can however also happen that the
			// result stream breaks halfway and ends with an error before
			// the end.
			it.err = err
			it.errIndex = it.nc
		}
		it.nc++
		if it.metadata == nil && it.RowIterator.Metadata != nil {
			it.metadata = it.RowIterator.Metadata
			// Initialize the checksum of the iterator by calculating the
			// checksum of the columns that are included in this result. This is
			// also used to detect the possible difference between two empty
			// result sets with a different set of columns.
			it.checksum, err = createMetadataChecksum(it.enc, it.buffer, it.metadata)
			if err != nil {
				return err
			}
		}
		if it.err != nil {
			return it.err
		}
		// Update the current checksum.
		it.checksum, err = updateChecksum(it.enc, it.buffer, it.checksum, row)
		return err
	})
	return row, err
}

// updateChecksum calculates the following checksum based on a current checksum
// and a new row.
func updateChecksum(enc *gob.Encoder, buffer *bytes.Buffer, currentChecksum *[32]byte, row *spanner.Row) (*[32]byte, error) {
	buffer.Reset()
	buffer.Write(currentChecksum[:])
	for i := 0; i < row.Size(); i++ {
		var v spanner.GenericColumnValue
		err := row.Column(i, &v)
		if err != nil {
			return nil, err
		}
		err = enc.Encode(v)
		if err != nil {
			return nil, err
		}
	}
	res := sha256.Sum256(buffer.Bytes())
	return &res, nil
}

// createMetadataChecksum calculates the checksum of the metadata of a result.
// Only the column names and types are included in the checksum. Any transaction
// metadata is not included.
func createMetadataChecksum(enc *gob.Encoder, buffer *bytes.Buffer, metadata *sppb.ResultSetMetadata) (*[32]byte, error) {
	buffer.Reset()
	for _, field := range metadata.RowType.Fields {
		err := enc.Encode(field)
		if err != nil {
			return nil, err
		}
	}
	res := sha256.Sum256(buffer.Bytes())
	return &res, nil
}

// retry implements retriableStatement.retry for queries. It will execute the
// query on a new Spanner transaction and iterate over the same number of rows
// as the initial attempt, and then compare the checksum of the initial and the
// retried iterator. It will also check if any error that was returned by the
// initial iterator was also returned by the new iterator, and that the errors
// were returned by the same row index.
func (it *checksumRowIterator) retry(ctx context.Context, tx *spanner.ReadWriteStmtBasedTransaction) error {
	buffer := &bytes.Buffer{}
	enc := gob.NewEncoder(buffer)
	retryIt := tx.Query(ctx, it.stmt)
	// If the original iterator had been stopped, we should also always stop the
	// new iterator.
	if it.stopped {
		defer retryIt.Stop()
	}
	// The underlying iterator will be replaced by the new one if the retry succeeds.
	replaceIt := func(err error) error {
		if it.RowIterator != nil {
			it.RowIterator.Stop()
			it.RowIterator = retryIt
		}
		it.metadata = retryIt.Metadata
		return err
	}
	// If the retry fails, we will not replace the underlying iterator and we should
	// stop the iterator that was used by the retry.
	failRetry := func(err error) error {
		retryIt.Stop()
		return err
	}
	// Iterate over the new result set as many times as we iterated over the initial
	// result set. The checksums of the two should be equal. Also, the new result set
	// should return any error on the same index as the original.
	var newChecksum *[32]byte
	var checksumErr error
	for n := int64(0); n < it.nc; n++ {
		row, err := retryIt.Next()
		if n == 0 && (err == nil || err == iterator.Done) {
			newChecksum, checksumErr = createMetadataChecksum(enc, buffer, retryIt.Metadata)
			if checksumErr != nil {
				return failRetry(checksumErr)
			}
		}
		if err != nil {
			if spanner.ErrCode(err) == codes.Aborted {
				// This fails this retry, but it will trigger a new retry of the
				// entire transaction.
				return failRetry(err)
			}
			if errorsEqualForRetry(err, it.err) && n == it.errIndex {
				// Check that the checksums are also equal.
				if !checksumsEqual(newChecksum, it.checksum) {
					return failRetry(ErrAbortedDueToConcurrentModification)
				}
				return replaceIt(nil)
			}
			return failRetry(ErrAbortedDueToConcurrentModification)
		}
		newChecksum, err = updateChecksum(enc, buffer, newChecksum, row)
		if err != nil {
			return failRetry(err)
		}
	}
	// Check if the initial attempt ended with an error and the current attempt
	// did not. This is normally an indication that the retry returned more
	// results than the initial attempt, and that the initial attempt returned
	// iterator.Done, but it could theoretically be any other error as well.
	if it.err != nil {
		return failRetry(ErrAbortedDueToConcurrentModification)
	}
	if !checksumsEqual(newChecksum, it.checksum) {
		return failRetry(ErrAbortedDueToConcurrentModification)
	}
	// Everything seems to be equal, replace the underlying iterator and return
	// a nil error.
	return replaceIt(nil)
}

func checksumsEqual(c1, c2 *[32]byte) bool {
	return (reflect.ValueOf(c1).IsNil() && reflect.ValueOf(c2).IsNil()) || *c1 == *c2
}

func (it *checksumRowIterator) Stop() {
	if !it.stopped {
		it.stopped = true
		it.RowIterator.Stop()
		it.RowIterator = nil
	}
}

func (it *checksumRowIterator) Metadata() *sppb.ResultSetMetadata {
	return it.metadata
}
