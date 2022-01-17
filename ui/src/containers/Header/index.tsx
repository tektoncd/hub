import React from 'react';
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
  TextListItem,
  Button,
  AlertVariant,
  Divider,
  DropdownItem,
  DropdownSeparator,
  Flex,
  FlexItem,
  Dropdown,
  DropdownToggle
} from '@patternfly/react-core';
import logo from '../../assets/logo/logo.png';
import { IconSize } from '@patternfly/react-icons';
import Search from '../../containers/Search';
import './Header.css';
import { scrollToTop } from '../../common/scrollToTop';
import Icon from '../../components/Icon';
import { Icons } from '../../common/icons';
import UserProfile from '../UserProfile';
import { useMst } from '../../store/root';
import { AUTH_BASE_URL, REDIRECT_URI } from '../../config/constants';
import { IProvider } from '../../store/provider';
import { titleCase } from '../../common/titlecase';
import AlertDisplay from '../../components/AlertDisplay';
import { ICatalog } from '../../store/catalog';

const Header: React.FC = observer(() => {
  const { user, providers, catalogs } = useMst();
  const history = useHistory();
  const [isModalOpen, setIsModalOpen] = React.useState(false);
  const [disable, setDisable] = React.useState(false);
  const [isResInstModallOpen, setResInstModalOpen] = React.useState(false);
  const [isCatInstModallOpen, setCatInstModalOpen] = React.useState(false);
  const [isOpen, setIsOpen] = React.useState(false);

  const dropdownItems = [
    <DropdownItem
      className="hub-header-contribute__dropdown"
      key="action-1"
      onClick={() => setResInstModalOpen(true)}
    >
      Add a resource
    </DropdownItem>,
    <DropdownSeparator key="separator-1" />,
    <DropdownItem
      className="hub-header-contribute__dropdown"
      key="action-2"
      onClick={() => setCatInstModalOpen(true)}
    >
      Add a catalog
    </DropdownItem>
  ];

  const headerTools = (
    <PageHeaderTools>
      <Flex>
        <FlexItem>
          <Flex>
            <FlexItem>
              <Search />
            </FlexItem>
            <FlexItem onClick={() => setIsModalOpen(true)} className="hub-header-search-hint">
              <Icon id={Icons.Help} size={IconSize.sm} label={'search-tips'} />
            </FlexItem>
          </Flex>
        </FlexItem>
        <FlexItem className="hub-header-contribute__margin">
          <Text component={TextVariants.h3}>
            <Dropdown
              onSelect={() => setIsOpen(!isOpen)}
              isPlain
              toggle={
                <DropdownToggle onToggle={() => setIsOpen(!isOpen)} className="hub-header-dropdown">
                  <Text component={TextVariants.h2}>
                    <span className="hub-header-contribute">Contribute</span>
                  </Text>
                </DropdownToggle>
              }
              isOpen={isOpen}
              dropdownItems={dropdownItems}
            />
          </Text>
        </FlexItem>
        <FlexItem>
          {user.isAuthenticated && user.refreshTokenInfo.expiresAt * 1000 > global.Date.now() ? (
            <UserProfile />
          ) : (
            <Text
              style={{ textDecoration: 'none' }}
              component={TextVariants.a}
              onClick={() => user.setIsAuthModalOpen(true)}
            >
              <span className="hub-header-login">
                <b>Login</b>
              </span>
            </Text>
          )}
        </FlexItem>
      </Flex>
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

      <Modal
        title="Instructions to add a new Resource in Hub"
        isOpen={isResInstModallOpen}
        onClose={() => setResInstModalOpen(false)}
        width={'50%'}
      >
        <Grid>
          <TextContent>
            <Text component={TextVariants.h3}>
              Add a new Resource to Hub by following the below steps
            </Text>
            <Text>
              Add a new resource to Hub by creating a pull request to available Catalogs where the
              new resource should follow the guildelines mentioned
              <Text
                component={TextVariants.a}
                target="_blank"
                href="https://github.com/tektoncd/catalog#catalog-structure"
              >
                {' '}
                here
              </Text>
            </Text>
            <Text>Available catalogs are:</Text>
            <TextList>
              {catalogs.values.map((catalog: ICatalog, index: number) => (
                <TextListItem key={index}>
                  <Text component={TextVariants.a} href={`${catalog.url}`} target="_blank">
                    {titleCase(catalog.name)}
                  </Text>
                </TextListItem>
              ))}
            </TextList>
            <Text component={TextVariants.blockquote}>
              Note: Newly added resource would be available on the Hub within 30 minutes once the PR
              gets merged.
            </Text>
          </TextContent>
        </Grid>
      </Modal>
      <Modal
        title="Instructions to add a new catalog in Hub"
        isOpen={isCatInstModallOpen}
        onClose={() => setCatInstModalOpen(false)}
        width={'50%'}
      >
        <Grid>
          <TextContent>
            <Text component={TextVariants.h3}>
              {' '}
              Add a new catalog to the list of available catalogs on Hub by modifying
              <Text
                component={TextVariants.a}
                target="_blank"
                href="https://github.com/tektoncd/hub/blob/main/config.yaml"
              >
                {' '}
                config.yaml .
              </Text>
              For more details please refer to the
              <Text
                component={TextVariants.a}
                target="_blank"
                href="https://github.com/tektoncd/hub/blob/main/docs/ADD_NEW_CATALOG.md"
              >
                {' '}
                document
              </Text>
              .
            </Text>

            <Text component={TextVariants.blockquote}>
              Note: If you are adding a new catalog to Hub then a config refresh needs to be done by
              the user who has config refresh scopes as per
              <Text
                component={TextVariants.a}
                href="https://github.com/tektoncd/hub/blob/main/config.yaml"
                target="_blank"
              >
                {' '}
                config.yaml
              </Text>
            </Text>
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
