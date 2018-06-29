import GearIcon from '@sourcegraph/icons/lib/Gear'
import PuzzleIcon from '@sourcegraph/icons/lib/Puzzle'
import UserIcon from '@sourcegraph/icons/lib/User'
import * as React from 'react'
import { Link, NavLink, RouteComponentProps } from 'react-router-dom'
import { platformEnabled } from '../../user/tags'
import { OrgAvatar } from '../OrgAvatar'
import { OrgAreaPageProps } from './OrgArea'

interface Props extends OrgAreaPageProps, RouteComponentProps<{}> {
    className: string
}

/**
 * Header for the organization area.
 */
export const OrgHeader: React.SFC<Props> = (props: Props) => (
    <div className={`org-header area-header ${props.className}`}>
        <div className={`${props.className}-inner`}>
            {props.org && (
                <>
                    <h2 className="org-header__title">
                        <OrgAvatar org={props.org.name} />{' '}
                        <span className="org-header__title-text">{props.org.displayName || props.org.name}</span>
                    </h2>
                    <div className="area-header__nav">
                        <div className="area-header__nav-links">
                            <NavLink
                                to={`${props.match.url}`}
                                exact={true}
                                className="btn area-header__nav-link"
                                activeClassName="area-header__nav-link--active"
                            >
                                Overview
                            </NavLink>
                            <NavLink
                                to={`${props.match.url}/members`}
                                exact={true}
                                className="btn area-header__nav-link"
                                activeClassName="area-header__nav-link--active"
                            >
                                <UserIcon className="icon-inline" /> Members
                            </NavLink>
                            {props.authenticatedUser &&
                                platformEnabled(props.authenticatedUser) && (
                                    <NavLink
                                        to={`${props.match.url}/extensions`}
                                        className="btn area-header__nav-link"
                                        activeClassName="area-header__nav-link--active"
                                    >
                                        <PuzzleIcon className="icon-inline" /> Extensions
                                    </NavLink>
                                )}
                            {props.org.viewerCanAdminister && (
                                <NavLink
                                    to={`${props.match.url}/settings`}
                                    className="btn area-header__nav-link"
                                    activeClassName="area-header__nav-link--active"
                                >
                                    <GearIcon className="icon-inline" /> Settings
                                </NavLink>
                            )}
                        </div>
                        {props.org.viewerPendingInvitation &&
                            props.org.viewerPendingInvitation.respondURL && (
                                <div className="area-header__nav-actions">
                                    <small className="area-header__nav-actions-label mr-sm-2">Join organization:</small>
                                    <Link
                                        to={props.org.viewerPendingInvitation.respondURL}
                                        className="btn btn-success btn-sm"
                                    >
                                        View invitation
                                    </Link>
                                </div>
                            )}
                    </div>
                </>
            )}
        </div>
    </div>
)
