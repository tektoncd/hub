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
  Grid
} from '@patternfly/react-core';
import logo from '../../assets/logo/logo.png';
import Search from '../../containers/Search';
import UserProfile from '../UserProfile';
import { useMst } from '../../store/root';
import './Header.css';

const Header: React.FC = observer(() => {
  const { user } = useMst();
  const history = useHistory();

  const headerTools = (
    <PageHeaderTools>
      <Grid>
        <GridItem span={11}>
          <Search />
        </GridItem>
      </Grid>
      {user.isAuthenticated ? (
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

  return (
    <React.Fragment>
      <PageHeader
        logo={<Brand src={logo} alt="Tekton Hub Logo" onClick={() => history.push('/')} />}
        headerTools={headerTools}
      />
    </React.Fragment>
  );
});
export default Header;
