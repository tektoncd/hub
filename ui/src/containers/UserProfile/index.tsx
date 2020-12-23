import React, { useState } from 'react';
import {
  Avatar,
  ClipboardCopy,
  ClipboardCopyVariant,
  Dropdown,
  DropdownItem,
  DropdownToggle,
  Modal
} from '@patternfly/react-core';
import imgAvatar from '../../assets/logo/imgAvatar.png';
import { useMst } from '../../store/root';
import './UserProfile.css';

const UserProfile: React.FC = () => {
  const { user } = useMst();

  const [isModalOpen, setIsModalOpen] = useState(false);
  const [isOpen, set] = useState(false);

  const onToggle = (isOpen: React.SetStateAction<boolean>) => set(isOpen);

  const dropdownItems = [
    <DropdownItem key="copyToken" onClick={() => setIsModalOpen(!isModalOpen)}>
      Copy Hub Token
    </DropdownItem>,
    <DropdownItem key="logout" onClick={user.logout}>
      Logout
    </DropdownItem>
  ];

  const userLogo: React.ReactNode = <Avatar className="hub-userlogo-size" src={imgAvatar} alt="" />;

  return (
    <React.Fragment>
      <Dropdown
        className="hub-userProfile"
        position="right"
        dropdownItems={dropdownItems}
        toggle={<DropdownToggle onToggle={onToggle}>{userLogo}</DropdownToggle>}
        isPlain
        isOpen={isOpen}
      ></Dropdown>
      <Modal
        className="hub-userProfile-modal"
        title="Copy Hub Token"
        isOpen={isModalOpen}
        onClose={() => setIsModalOpen(!isModalOpen)}
      >
        <hr />
        <div>
          <ClipboardCopy isReadOnly variant={ClipboardCopyVariant.expansion}>
            {user.accessTokenInfo.token}
          </ClipboardCopy>
        </div>
      </Modal>
    </React.Fragment>
  );
};
export default UserProfile;
