import React, { useState } from 'react';
import { observer } from 'mobx-react';
import { useHistory, useParams } from 'react-router-dom';
import { IconSize } from '@patternfly/react-icons';
import { Alert, AlertActionCloseButton, AlertGroup } from '@patternfly/react-core';
import { useMst } from '../../store/root';
import { assert } from '../../store/utils';
import { Icons } from '../../common/icons';
import Icon from '../../components/Icon';
import './Rating.css';

const Rating: React.FC = observer(() => {
  const { name } = useParams();

  const [star, setStar] = useState([false, false, false, false, false]);
  const [ratingError, setRatingError] = useState('');

  const highlightStar = (value: number) => {
    const checkedStatus = [false, false, false, false, false];
    switch (value) {
      case 5:
        checkedStatus[4] = true;
        break;
      case 4:
        checkedStatus[3] = true;
        break;
      case 3:
        checkedStatus[2] = true;
        break;
      case 2:
        checkedStatus[1] = true;
        break;
      case 1:
        checkedStatus[0] = true;
        break;
      default:
    }
    setStar(checkedStatus);
  };

  const starList = (
    <ul className="hub-rate-area">
      <input readOnly type="radio" id="5" name="rating" value="5" checked={star[4]} />
      <label htmlFor="5">
        <Icon id={Icons.Star} size={IconSize.sm} label="5" />
      </label>
      <input readOnly type="radio" id="4" name="rating" value="4" checked={star[3]} />
      <label htmlFor="4">
        <Icon id={Icons.Star} size={IconSize.sm} label="4" />
      </label>
      <input readOnly type="radio" id="3" name="rating" value="3" checked={star[2]} />
      <label htmlFor="3">
        <Icon id={Icons.Star} size={IconSize.sm} label="3" />
      </label>
      <input readOnly type="radio" id="2" name="rating" value="2" checked={star[1]} />
      <label htmlFor="2">
        <Icon id={Icons.Star} size={IconSize.sm} label="2" />
      </label>
      <input readOnly type="radio" id="1" name="rating" value="1" checked={star[0]} />
      <label htmlFor="1">
        <Icon id={Icons.Star} size={IconSize.sm} label="1" />
      </label>
    </ul>
  );

  const { user, resources } = useMst();
  const history = useHistory();
  const rateResource = (event: React.MouseEvent<HTMLFormElement, MouseEvent>) => {
    const target = event.target as HTMLInputElement;
    const rating = target.value;

    if (!user.isAuthenticated) {
      history.push('/login');
    } else {
      if (rating) {
        const resource = resources.resources.get(name);
        assert(resource);
        user.setRating(resource.id, Number(rating));
        if (user.ratingErr === '') {
          highlightStar(Number(rating));
        } else {
          setRatingError(user.ratingErr);
        }
      }
    }
  };

  if (!user.isLoading) {
    highlightStar(user.userRating);
    user.setLoading(true);
  }

  return (
    <div className="hub-details-rating">
      <form onClick={rateResource}>{starList}</form>
      {ratingError ? (
        <AlertGroup isToast>
          <Alert
            isLiveRegion
            variant="info"
            title={ratingError}
            actionClose={
              <AlertActionCloseButton
                onClose={() => {
                  window.location.reload();
                }}
              />
            }
          ></Alert>
        </AlertGroup>
      ) : null}
    </div>
  );
});
export default Rating;
