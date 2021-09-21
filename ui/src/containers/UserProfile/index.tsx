import React, { useState, useEffect, useCallback } from 'react';
import {
  Avatar,
  ClipboardCopy,
  ClipboardCopyVariant,
  Dropdown,
  DropdownItem,
  DropdownToggle,
  Modal,
  DropdownSeparator
} from '@patternfly/react-core';
import { observer } from 'mobx-react';
import { useMst } from '../../store/root';
import './UserProfile.css';

const UserProfile: React.FC = observer(() => {
  const { user } = useMst();

  const [refreshId, setRefreshId] = useState<number>(0);
  const [accessId, setAccessId] = useState<number>(0);

  useEffect(() => {
    user.getProfile();
  }, [user.profile.userName]);

  const triggerInterval = useCallback(() => {
    const accessTokenInterval = user.accessTokenInfo.expiresAt * 1000 - new Date().getTime();
    const refreshTokenInterval = user.refreshTokenInfo.expiresAt * 1000 - new Date().getTime();

    // The condition checks the maximum delay for setInterval
    if (refreshTokenInterval < Math.pow(2, 31) - 1) {
      // To get a new refresh token
      // Update the refresh token before 20 seconds of current refresh token's expiry time
      const tempRefreshId = window.setInterval(() => {
        if (!user.isAuthenticated) {
          clearInterval(tempRefreshId);
        }
        user.updateRefreshToken();
      }, refreshTokenInterval - 20000);
      setRefreshId(tempRefreshId);
    }

    // The condition checks the maximum delay for setInterval
    if (accessTokenInterval < Math.pow(2, 31) - 1) {
      // To get a new access token
      // Update the access token before 10 seconds of current access token's expiry time
      const tempAccessId = window.setInterval(() => {
        if (!user.isAuthenticated) {
          clearInterval(tempAccessId);
        }
        user.updateAccessToken();
      }, accessTokenInterval - 10000);
      setAccessId(tempAccessId);
    }
  }, [user]);

  useEffect(() => {
    triggerInterval();
  }, [triggerInterval]);

  const [isModalOpen, setIsModalOpen] = useState(false);
  const [isOpen, set] = useState(false);

  const hubLogout = () => {
    user.logout();
    localStorage.clear();
    clearInterval(refreshId);
    clearInterval(accessId);
  };

  const onToggle = (isOpen: React.SetStateAction<boolean>) => set(isOpen);

  const dropdownItems = [
    <DropdownItem key="githubId">Hi {user.profile.userName}</DropdownItem>,
    <DropdownSeparator key="separator" />,
    <DropdownItem key="copyToken" onClick={() => setIsModalOpen(!isModalOpen)}>
      Copy Hub Token
    </DropdownItem>,
    <DropdownItem key="logout" onClick={hubLogout}>
      Logout
    </DropdownItem>
  ];

  const userLogo: React.ReactNode = (
    <Avatar className="hub-userlogo-size" src={user.profile.avatarUrl} alt="" />
  );

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
});
export default UserProfile;
