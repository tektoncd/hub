import React from 'react';
import { observer } from 'mobx-react';
import { useNavigate } from 'react-router-dom';
import {
  PageHeader,
  Brand,
  PageHeaderTools,
  Text,
  TextVariants,
  GridItem,
  Grid,
  Modal,
  ModalVariant,
  TextContent,
  TextList,
  TextListItem,
  Button,
  AlertVariant,
  Divider
} from '@patternfly/react-core';
import logo from '../../assets/logo/logo.png';
import { IconSize } from '@patternfly/react-icons';
import Search from '../../containers/Search';
import { scrollToTop } from '../../common/scrollToTop';
import Icon from '../../components/Icon';
import { Icons } from '../../common/icons';
import UserProfile from '../UserProfile';
import { useMst } from '../../store/root';
import { AUTH_BASE_URL, REDIRECT_URI } from '../../config/constants';
import { IProvider } from '../../store/provider';
import { titleCase } from '../../common/titlecase';
import AlertDisplay from '../../components/AlertDisplay';
import './Header.css';
import { apiDownError } from '../../common/errors';

const Header: React.FC = observer(() => {
  const { user, resources, providers } = useMst();
  const history = useNavigate();
  const [isModalOpen, setIsModalOpen] = React.useState(false);
  const [disable, setDisable] = React.useState(false);

  const headerTools = (
    <PageHeaderTools>
      <Grid>
        <GridItem span={AUTH_BASE_URL !== '' ? 10 : 11}>
          <Search />
        </GridItem>
        <GridItem span={1} onClick={() => setIsModalOpen(true)} className="hub-header-search-hint">
          <Icon id={Icons.Help} size={IconSize.sm} label={'search-tips'} />
        </GridItem>
      </Grid>
      {user.isAuthenticated && user.refreshTokenInfo.expiresAt * 1000 > global.Date.now() ? (
        <UserProfile />
      ) : AUTH_BASE_URL !== '' ? (
        <Text
          style={{ textDecoration: 'none' }}
          component={TextVariants.a}
          onClick={() => (resources.err !== apiDownError ? user.setIsAuthModalOpen(true) : null)}
        >
          <span className="hub-header-login">
            <b>Login</b>
          </span>
        </Text>
      ) : null}
    </PageHeaderTools>
  );

  const homePage = () => {
    if (!window.location.search) history('/');
    scrollToTop();
  };

  return (
    <React.Fragment>
      <PageHeader
        logo={<Brand src={logo} alt="Tekton Hub Logo" onClick={homePage} />}
        headerTools={headerTools}
      />

      <Modal
        variant={ModalVariant.small}
        title="Search tips:"
        isOpen={isModalOpen}
        onClose={() => setIsModalOpen(false)}
      >
        <Grid>
          <TextContent>
            <TextList>
              <TextListItem>Press `/` to quickly focus on search.</TextListItem>
              <TextListItem>Search resources by name, displayName, and tags.</TextListItem>
              <TextListItem>
                Filter resources by tags using the qualifier like `tags:tagA,tagB`
              </TextListItem>
            </TextList>
          </TextContent>
        </Grid>
      </Modal>

      <Modal
        variant={ModalVariant.small}
        title={`Tekton Hub`}
        isOpen={user.isAuthModalOpen}
        onClose={() => user.setIsAuthModalOpen(false)}
        className="hub-header-login__modal"
        aria-label="login"
      >
        <TextContent>
          <Divider />
          <Text component={TextVariants.h6}>Sign In With</Text>
        </TextContent>
        <Grid>
          {providers.values.map((provider: IProvider) => (
            <GridItem key={provider.name} offset={1} span={10} className="hub-header-sigin-button">
              <Button
                variant="tertiary"
                component="a"
                isDisabled={disable}
                onClick={() => setDisable(true)}
                isBlock
                href={`${AUTH_BASE_URL}/auth/${provider.name}?redirect_uri=${REDIRECT_URI}`}
              >
                <span className="hub-header-sigin-button__icon ">
                  <Icon id={provider.name as Icons} size={IconSize.sm} label={provider.name} />
                </span>
                {titleCase(provider.name)}
              </Button>
            </GridItem>
          ))}
        </Grid>
      </Modal>
      {user.authErr.serverMessage ? (
        <AlertDisplay message={user.authErr} alertVariant={AlertVariant.danger} />
      ) : null}
    </React.Fragment>
  );
});
export default Header;
