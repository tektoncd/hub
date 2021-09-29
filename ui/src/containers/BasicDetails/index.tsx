import React, { useState } from 'react';
import {
  Card,
  CardHeader,
  Grid,
  GridItem,
  Text,
  TextVariants,
  Badge,
  CardActions,
  Button,
  Dropdown,
  DropdownToggle,
  Modal,
  TextContent,
  ClipboardCopy,
  ClipboardCopyVariant,
  DropdownItem,
  Spinner,
  List,
  ListItem,
  ListVariant
} from '@patternfly/react-core';

import { IconSize } from '@patternfly/react-icons';
import { useObserver } from 'mobx-react';
import { useParams } from 'react-router-dom';
import { useMst } from '../../store/root';
import { IResource, IResourceStore } from '../../store/resource';
import { ITag } from '../../store/tag';
import { Icons } from '../../common/icons';
import Icon from '../../components/Icon';
import TooltipDisplay from '../../components/TooltipDisplay';
import Rating from '../Rating';
import { titleCase } from '../../common/titlecase';
import { assert } from '../../store/utils';
import { ICatalogStore } from '../../store/catalog';
import './BasicDetails.css';

const BasicDetails: React.FC = () => {
  const { resources, catalogs }: { resources: IResourceStore; catalogs: ICatalogStore } = useMst();
  const { catalog, kind, name }: { catalog: string; kind: string; name: string } = useParams();

  const resourceKey = `${catalog}/${titleCase(kind)}/${name}`;
  const resource: IResource | undefined = resources.resources.get(resourceKey);
  assert(resource);

  const catalogId = resource.catalog.id;
  let catalogProvider = 'Github';

  for (let catalogItem = 0; catalogItem < catalogs.values.length; catalogItem++) {
    if (catalogs.values[catalogItem].id === catalogId) {
      catalogProvider = catalogs.values[catalogItem].provider;
      break;
    }
  }

  const dropdownItems = resource.versions.map((value) => (
    <DropdownItem
      id={String(value.id)}
      key={value.id}
      onClick={(e) => resources.setDisplayVersion(resourceKey, e.currentTarget.id)}
    >
      {value.version === resource.latestVersion.version
        ? `${value.version} (latest)`
        : value.version}
    </DropdownItem>
  ));

  const [isOpen, set] = useState(false);
  const onToggle = (isOpen: React.SetStateAction<boolean>) => set(isOpen);
  const onSelect = () => set(!isOpen);

  const [isModalOpen, setIsModalOpen] = useState(false);
  const onModalToggle = () => setIsModalOpen(!isModalOpen);

  return useObserver(() =>
    resource.versions.length === 0 ? (
      <Spinner />
    ) : (
      <Card className="hub-header-card">
        <Grid className="hub-header-card__margin">
          <GridItem span={1}>
            <Grid>
              <GridItem offset={8}>
                <div className="hub-details-kind-icon">
                  <Icon id={resource.kind.icon} size={IconSize.xl} label={resource.kind.name} />
                </div>
              </GridItem>
            </Grid>
          </GridItem>
          <GridItem span={10}>
            <CardHeader>
              <TextContent className="hub-details-card-body">
                <Grid className="hub-details-title">
                  <GridItem span={11}>
                    <List variant={ListVariant.inline} style={{ listStyleType: 'none' }}>
                      <ListItem>
                        <Text className="hub-details-resource-name">{resource.resourceName}</Text>
                      </ListItem>
                      <ListItem>
                        <Icon
                          id={resource.catalog.icon}
                          size={IconSize.lg}
                          label={resource.catalog.name}
                        />
                      </ListItem>
                      <ListItem>
                        {resource.displayVersion.deprecated ? (
                          <TooltipDisplay
                            name={'This version has been deprecated'}
                            id={Icons.WarningTriangle}
                            size={IconSize.lg}
                          />
                        ) : null}
                      </ListItem>
                    </List>
                  </GridItem>
                </Grid>
                <a href={resource.webURL} target="_" className="hub-details-hyperlink">
                  <List variant={ListVariant.inline} style={{ listStyleType: 'none' }}>
                    <ListItem>
                      <Icon id={Icons.Github} size={IconSize.md} label="Github" />
                    </ListItem>
                    <ListItem>
                      <Text className="hub-details-github">
                        Open {resource.name} in {titleCase(catalogProvider)}
                      </Text>
                    </ListItem>
                  </List>
                </a>
                <Grid>
                  <GridItem span={10} className="hub-details-description">
                    <div className="line">{resource.summary}</div>
                    <div>{resource.detailDescription}</div>
                  </GridItem>
                  <GridItem>
                    {resource.tags.map((tag: ITag) => (
                      <Badge key={`badge-${tag.id}`} className="hub-tags">
                        {tag.name}
                      </Badge>
                    ))}
                  </GridItem>
                </Grid>
              </TextContent>
              <CardActions className="hub-details-card-action">
                <Grid>
                  <GridItem offset={2} span={1}>
                    <Grid>
                      <GridItem offset={3}>
                        <TooltipDisplay id={Icons.Star} size={IconSize.sm} name="Average Rating" />
                      </GridItem>
                    </Grid>
                  </GridItem>
                  <GridItem span={1}>
                    <Text> {resource.rating}</Text>
                  </GridItem>
                  <GridItem className="hub-details-rating__margin">
                    <Rating />
                  </GridItem>
                  <GridItem className="hub-details-rating__margin">
                    <Button
                      variant="primary"
                      className="hub-details-button"
                      onClick={onModalToggle}
                    >
                      Install
                    </Button>
                  </GridItem>
                  <GridItem className="hub-details-rating__margin">
                    <Dropdown
                      toggle={
                        <DropdownToggle onToggle={onToggle} className="hub-details-dropdown-item">
                          {resource.displayVersion.id === resource.latestVersion.id
                            ? `${resource.displayVersion.version} (latest)`
                            : `${resource.displayVersion.version}`}
                        </DropdownToggle>
                      }
                      dropdownItems={dropdownItems}
                      onSelect={onSelect}
                      isOpen={isOpen}
                    />
                  </GridItem>
                </Grid>
              </CardActions>
            </CardHeader>
          </GridItem>

          <Modal width={'60%'} title={resource.name} isOpen={isModalOpen} onClose={onModalToggle}>
            <hr />
            <div>
              <TextContent>
                <Text component={TextVariants.h2}>Install using kubectl</Text>
                <Text> {resource.kind.name} </Text>
                <ClipboardCopy isReadOnly variant={ClipboardCopyVariant.expansion}>
                  {resource.installCommand}
                </ClipboardCopy>
              </TextContent>
              <br />
              <TextContent>
                <Text component={TextVariants.h2}>Install using tkn</Text>
                <Text> {resource.kind.name} </Text>
                <ClipboardCopy isReadOnly variant={ClipboardCopyVariant.expansion}>
                  {resource.tknInstallCommand}
                </ClipboardCopy>
              </TextContent>
            </div>
          </Modal>
        </Grid>
      </Card>
    )
  );
};

export default BasicDetails;
