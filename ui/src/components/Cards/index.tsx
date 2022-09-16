import React from 'react';
import {
  Card,
  CardHeader,
  CardActions,
  CardTitle,
  CardBody,
  TextContent,
  Badge,
  CardFooter,
  GalleryItem
} from '@patternfly/react-core';
import { Link } from 'react-router-dom';
import { IconSize } from '@patternfly/react-icons';
import { IResource } from '../../store/resource';
import { ITag } from '../../store/tag';
import Icon from '../Icon';
import { Icons } from '../../common/icons';
import TooltipDisplay from '../TooltipDisplay';
import './Cards.css';

interface Props {
  items: IResource[];
}

const Cards: React.FC<Props> = (resources) => {
  return (
    <React.Fragment>
      {resources.items.map((resource: IResource, r: number) => (
        <GalleryItem key={r}>
          <Link
            to={{
              pathname: `${resource.catalog.name.toLowerCase()}/${resource.kind.name.toLowerCase()}/${resource.name.toLowerCase()}`
            }}
            className="hub-card-link"
          >
            <Card className="hub-resource-card" data-test={resource.name}>
              <CardHeader>
                <span className="hub-kind-icon">
                  <TooltipDisplay name={resource.kind.name} id={resource.kind.icon} />
                </span>
                <span className="hub-catalog-icon">
                  <TooltipDisplay name={resource.catalog.name} id={resource.catalog.icon} />
                </span>
                <CardActions>
                  <span className="hub-rating-icon">
                    <Icon id={Icons.Star} size={IconSize.sm} label={String(resource.rating)} />
                  </span>
                  <TextContent className="hub-rating">{resource.rating}</TextContent>
                </CardActions>
              </CardHeader>

              <CardTitle>
                <span className="hub-resource-name">{resource.resourceName}</span>
                <span className="hub-resource-version">v{resource.latestVersion.version}</span>
              </CardTitle>

              <CardBody className="hub-resource-body fade">
                {resource.latestVersion.description}
              </CardBody>

              <CardFooter>
                <TextContent className="hub-resource-updatedAt">
                  Updated {resource.latestVersion.updatedAt.fromNow()}
                </TextContent>

                <div className="hub-tags-container">
                  {resource.tags.slice(0, 3).map((tag: ITag) => (
                    <Badge className="hub-tags" key={`badge-${tag.id}`}>
                      {tag.name}
                    </Badge>
                  ))}
                </div>
              </CardFooter>
            </Card>
          </Link>
        </GalleryItem>
      ))}
    </React.Fragment>
  );
};

export default Cards;
