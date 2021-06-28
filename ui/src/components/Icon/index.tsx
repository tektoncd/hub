import React from 'react';
import {
  CatIcon,
  CertificateIcon,
  BuildIcon,
  DomainIcon,
  UserIcon,
  IconSize,
  StarIcon,
  GithubIcon,
  WarningTriangleIcon,
  CatalogIcon,
  HelpIcon
} from '@patternfly/react-icons';
import { Icons } from '../../common/icons';
import './Icon.css';

interface Props {
  id: Icons;
  size: IconSize | keyof typeof IconSize;
  label: string;
}

const Icon: React.FC<Props> = (props: Props) => {
  const { id, size, label } = props;
  switch (id) {
    case Icons.Unknown:
      return <div></div>;
    case Icons.Cat:
      return <CatIcon size={size} className="hub-icon" label={label} />;
    case Icons.Certificate:
      return <CertificateIcon size={size} className="hub-icon" label={label} />;
    case Icons.User:
      return <UserIcon size={size} className="hub-icon" label={label} />;
    case Icons.Build:
      return <BuildIcon size={size} className="hub-icon" label={label} />;
    case Icons.Domain:
      return <DomainIcon size={size} className="hub-icon" label={label} />;
    case Icons.Star:
      return <StarIcon size={size} label={label} />;
    case Icons.Github:
      return <GithubIcon size={size} className="hub-icon" label={label} />;
    case Icons.WarningTriangle:
      return <WarningTriangleIcon size={size} className="hub-icon" label={label} />;
    case Icons.Catalog:
      return <CatalogIcon size={size} className="hub-icon" label={label} />;
    case Icons.Help:
      return <HelpIcon size={size} label={label} />;
  }
};

export default Icon;
