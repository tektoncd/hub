import React from 'react';
import { IconSize } from '@patternfly/react-icons';
import { Tooltip } from '@patternfly/react-core';
import Icon from '../Icon';
import { titleCase } from '../../common/titlecase';
import { Icons } from '../../common/icons';

interface Props {
  id: Icons;
  size?: IconSize;
  name: string;
}

const TooltipDisplay: React.FC<Props> = (props: Props) => {
  return (
    <Tooltip content={<b>{titleCase(props.name)}</b>}>
      <Icon id={props.id} size={props.size || IconSize.sm} label={props.name} />
    </Tooltip>
  );
};

export default TooltipDisplay;
