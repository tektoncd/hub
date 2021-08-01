import React from 'react';
import { Link } from 'react-router-dom';
import { observer } from 'mobx-react';
import { useHistory } from 'react-router-dom';
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
  TextListItem
} from '@patternfly/react-core';
import logo from '../../assets/logo/logo.png';
import { IconSize } from '@patternfly/react-icons';
import Search from '../../containers/Search';
import UserProfile from '../UserProfile';
import { useMst } from '../../store/root';
import './Header.css';
import { scrollToTop } from '../../common/scrollToTop';
import Icon from '../../components/Icon';
import { Icons } from '../../common/icons';

const Header: React.FC = observer(() => {
  const { user } = useMst();
  const history = useHistory();
  const [isModalOpen, setIsModalOpen] = React.useState(false);

  const headerTools = (
    <PageHeaderTools>
      <Grid>
        <GridItem span={10}>
          <Search />
        </GridItem>
        <GridItem span={1} onClick={() => setIsModalOpen(true)} className="header-search-hint">
          <Icon id={Icons.Help} size={IconSize.sm} label={'search-tips'} />
        </GridItem>
      </Grid>
      {user.isAuthenticated && user.refreshTokenInfo.expiresAt * 1000 > global.Date.now() ? (
        <UserProfile />
      ) : (
        <Text component={TextVariants.h3}>
          <Link to="/login" style={{ textDecoration: 'none' }}>
            <span className="hub-header-login">Login</span>
          </Link>
        </Text>
      )}
    </PageHeaderTools>
  );

  const homePage = () => {
    if (!window.location.search) history.push('/');
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
    </React.Fragment>
  );
});
export default Header;
