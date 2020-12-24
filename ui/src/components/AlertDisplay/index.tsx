import React from 'react';
import { Alert, AlertActionCloseButton, AlertGroup, AlertVariant } from '@patternfly/react-core';
import { IError } from '../../store/auth';
import './AlertDisplay.css';

interface Error {
  message: IError;
  alertVariant: AlertVariant | keyof typeof AlertVariant;
}

const AlertDisplay: React.FC<Error> = (err: Error) => {
  return (
    <AlertGroup isToast className="hub-alert">
      <Alert
        isLiveRegion
        variant={err.alertVariant}
        title={err.message.customMessage}
        actionClose={<AlertActionCloseButton onClose={() => window.location.reload()} />}
      >
        {err.message.serverMessage}
      </Alert>
    </AlertGroup>
  );
};

export default AlertDisplay;
