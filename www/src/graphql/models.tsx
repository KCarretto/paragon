import gql from 'graphql-tag';
export type Maybe<T> = T | null;
/** All built-in and custom scalars, mapped to their actual values */
export type Scalars = {
  ID: string,
  String: string,
  Boolean: boolean,
  Int: number,
  Float: number,
  Time: any,
};



export type ActivateServiceRequest = {
  id: Scalars['ID'],
};

export type ActivateUserRequest = {
  id: Scalars['ID'],
};

export type AddCredentialForTargetRequest = {
  id: Scalars['ID'],
  principal: Scalars['String'],
  secret: Scalars['String'],
  kind?: Maybe<Scalars['String']>,
};

export type AddCredentialForTargetsRequest = {
  ids?: Maybe<Array<Scalars['ID']>>,
  principal: Scalars['String'],
  secret: Scalars['String'],
  kind?: Maybe<Scalars['String']>,
};

export type ApplyTagRequest = {
  tagID: Scalars['ID'],
  entID: Scalars['ID'],
};

export type ApplyTagToTargetsRequest = {
  tagID: Scalars['ID'],
  targets?: Maybe<Array<Scalars['ID']>>,
};

export type ChangeNameRequest = {
  name: Scalars['String'],
};

export type ClaimTasksRequest = {
  machineUUID?: Maybe<Scalars['String']>,
  primaryIP?: Maybe<Scalars['String']>,
  hostname?: Maybe<Scalars['String']>,
  primaryMAC?: Maybe<Scalars['String']>,
  sessionID?: Maybe<Scalars['String']>,
};

export type CreateJobRequest = {
  name: Scalars['String'],
  content: Scalars['String'],
  sessionID?: Maybe<Scalars['String']>,
  stage?: Maybe<Scalars['Boolean']>,
  targets?: Maybe<Array<Scalars['ID']>>,
  tags?: Maybe<Array<Scalars['ID']>>,
  prev?: Maybe<Scalars['ID']>,
};

export type CreateLinkRequest = {
  alias: Scalars['String'],
  expirationTime?: Maybe<Scalars['Time']>,
  clicks?: Maybe<Scalars['Int']>,
  file: Scalars['ID'],
};

export type CreateTagRequest = {
  name: Scalars['String'],
};

export type CreateTargetRequest = {
  name: Scalars['String'],
  primaryIP: Scalars['String'],
  tags?: Maybe<Array<Scalars['ID']>>,
};

export type Credential = {
   __typename?: 'Credential',
  id: Scalars['ID'],
  principal?: Maybe<Scalars['String']>,
  secret?: Maybe<Scalars['String']>,
  kind?: Maybe<Scalars['String']>,
  fails?: Maybe<Scalars['Int']>,
};

export type DeactivateServiceRequest = {
  id: Scalars['ID'],
};

export type DeactivateUserRequest = {
  id: Scalars['ID'],
};

export type DeleteCredentialRequest = {
  id: Scalars['ID'],
};

export type DeleteTargetRequest = {
  id: Scalars['ID'],
};

export type Event = {
   __typename?: 'Event',
  id: Scalars['ID'],
  creationTime?: Maybe<Scalars['Time']>,
  kind?: Maybe<Scalars['String']>,
  job?: Maybe<Job>,
  file?: Maybe<File>,
  credential?: Maybe<Credential>,
  link?: Maybe<Link>,
  tag?: Maybe<Tag>,
  target?: Maybe<Target>,
  task?: Maybe<Task>,
  user?: Maybe<User>,
  service?: Maybe<Service>,
  event?: Maybe<Event>,
  likers?: Maybe<Array<Maybe<User>>>,
  owner?: Maybe<User>,
  svcOwner?: Maybe<Service>,
};


export type EventLikersArgs = {
  input?: Maybe<Filter>
};

export type FailCredentialRequest = {
  id: Scalars['ID'],
};

export type File = {
   __typename?: 'File',
  id: Scalars['ID'],
  name?: Maybe<Scalars['String']>,
  creationTime?: Maybe<Scalars['Time']>,
  lastModifiedTime?: Maybe<Scalars['Time']>,
  size?: Maybe<Scalars['Int']>,
  hash?: Maybe<Scalars['String']>,
  contentType?: Maybe<Scalars['String']>,
  links?: Maybe<Array<Maybe<Link>>>,
};


export type FileLinksArgs = {
  input?: Maybe<Filter>
};

export type Filter = {
  offset?: Maybe<Scalars['Int']>,
  limit?: Maybe<Scalars['Int']>,
  search?: Maybe<Scalars['String']>,
};

export type Job = {
   __typename?: 'Job',
  id: Scalars['ID'],
  name?: Maybe<Scalars['String']>,
  creationTime?: Maybe<Scalars['Time']>,
  content?: Maybe<Scalars['String']>,
  staged?: Maybe<Scalars['Boolean']>,
  tasks?: Maybe<Array<Maybe<Task>>>,
  tags?: Maybe<Array<Maybe<Tag>>>,
  next?: Maybe<Job>,
  prev?: Maybe<Job>,
  owner?: Maybe<User>,
};


export type JobTasksArgs = {
  input?: Maybe<Filter>
};


export type JobTagsArgs = {
  input?: Maybe<Filter>
};

export type LikeEventRequest = {
  id: Scalars['ID'],
};

export type Link = {
   __typename?: 'Link',
  id: Scalars['ID'],
  alias?: Maybe<Scalars['String']>,
  expirationTime?: Maybe<Scalars['Time']>,
  clicks?: Maybe<Scalars['Int']>,
  file?: Maybe<File>,
};

export type MakeAdminRequest = {
  id: Scalars['ID'],
};

export type Mutation = {
   __typename?: 'Mutation',
  /** Credential Mutations */
  failCredential: Credential,
  deleteCredential: Scalars['Boolean'],
  /** Job Mutations */
  createJob: Job,
  queueJob: Job,
  /** Tag Mutations */
  createTag: Tag,
  applyTagToTask: Task,
  applyTagToTargets?: Maybe<Array<Target>>,
  applyTagToJob: Job,
  removeTagFromTask: Task,
  removeTagFromTarget: Target,
  removeTagFromJob: Job,
  /** Target Mutations */
  createTarget: Target,
  setTargetFields: Target,
  deleteTarget: Scalars['Boolean'],
  addCredentialForTarget: Target,
  addCredentialForTargets?: Maybe<Array<Target>>,
  /** Task Mutations */
  claimTasks?: Maybe<Array<Task>>,
  claimTask?: Maybe<Task>,
  submitTaskResult: Task,
  /** Link Mutations */
  createLink: Link,
  setLinkFields: Link,
  /** User Mutations */
  activateUser: User,
  deactivateUser: User,
  makeAdmin: User,
  removeAdmin: User,
  changeName: User,
  /** Service Mutations */
  activateService: Service,
  deactivateService: Service,
  setServiceConfig: Service,
  /** Event Mutations */
  likeEvent: Event,
};


export type MutationFailCredentialArgs = {
  input?: Maybe<FailCredentialRequest>
};


export type MutationDeleteCredentialArgs = {
  input?: Maybe<DeleteCredentialRequest>
};


export type MutationCreateJobArgs = {
  input?: Maybe<CreateJobRequest>
};


export type MutationQueueJobArgs = {
  input?: Maybe<QueueJobRequest>
};


export type MutationCreateTagArgs = {
  input?: Maybe<CreateTagRequest>
};


export type MutationApplyTagToTaskArgs = {
  input?: Maybe<ApplyTagRequest>
};


export type MutationApplyTagToTargetsArgs = {
  input?: Maybe<ApplyTagToTargetsRequest>
};


export type MutationApplyTagToJobArgs = {
  input?: Maybe<ApplyTagRequest>
};


export type MutationRemoveTagFromTaskArgs = {
  input?: Maybe<RemoveTagRequest>
};


export type MutationRemoveTagFromTargetArgs = {
  input?: Maybe<RemoveTagRequest>
};


export type MutationRemoveTagFromJobArgs = {
  input?: Maybe<RemoveTagRequest>
};


export type MutationCreateTargetArgs = {
  input?: Maybe<CreateTargetRequest>
};


export type MutationSetTargetFieldsArgs = {
  input?: Maybe<SetTargetFieldsRequest>
};


export type MutationDeleteTargetArgs = {
  input?: Maybe<DeleteTargetRequest>
};


export type MutationAddCredentialForTargetArgs = {
  input?: Maybe<AddCredentialForTargetRequest>
};


export type MutationAddCredentialForTargetsArgs = {
  input?: Maybe<AddCredentialForTargetsRequest>
};


export type MutationClaimTasksArgs = {
  input?: Maybe<ClaimTasksRequest>
};


export type MutationClaimTaskArgs = {
  id: Scalars['ID']
};


export type MutationSubmitTaskResultArgs = {
  input?: Maybe<SubmitTaskResultRequest>
};


export type MutationCreateLinkArgs = {
  input?: Maybe<CreateLinkRequest>
};


export type MutationSetLinkFieldsArgs = {
  input?: Maybe<SetLinkFieldsRequest>
};


export type MutationActivateUserArgs = {
  input?: Maybe<ActivateUserRequest>
};


export type MutationDeactivateUserArgs = {
  input?: Maybe<DeactivateUserRequest>
};


export type MutationMakeAdminArgs = {
  input?: Maybe<MakeAdminRequest>
};


export type MutationRemoveAdminArgs = {
  input?: Maybe<RemoveAdminRequest>
};


export type MutationChangeNameArgs = {
  input?: Maybe<ChangeNameRequest>
};


export type MutationActivateServiceArgs = {
  input?: Maybe<ActivateServiceRequest>
};


export type MutationDeactivateServiceArgs = {
  input?: Maybe<DeactivateServiceRequest>
};


export type MutationSetServiceConfigArgs = {
  input?: Maybe<SetServiceConfigRequest>
};


export type MutationLikeEventArgs = {
  input?: Maybe<LikeEventRequest>
};

export type Query = {
   __typename?: 'Query',
  link?: Maybe<Link>,
  links?: Maybe<Array<Maybe<Link>>>,
  file?: Maybe<File>,
  files?: Maybe<Array<Maybe<File>>>,
  credential?: Maybe<Credential>,
  credentials?: Maybe<Array<Maybe<Credential>>>,
  job?: Maybe<Job>,
  jobs?: Maybe<Array<Maybe<Job>>>,
  tag?: Maybe<Tag>,
  tags?: Maybe<Array<Maybe<Tag>>>,
  target?: Maybe<Target>,
  targets?: Maybe<Array<Maybe<Target>>>,
  task?: Maybe<Task>,
  tasks?: Maybe<Array<Maybe<Task>>>,
  user?: Maybe<User>,
  me?: Maybe<User>,
  users?: Maybe<Array<Maybe<User>>>,
  service?: Maybe<Service>,
  services?: Maybe<Array<Maybe<Service>>>,
  event?: Maybe<Event>,
  events?: Maybe<Array<Maybe<Event>>>,
};


export type QueryLinkArgs = {
  id: Scalars['ID']
};


export type QueryLinksArgs = {
  input?: Maybe<Filter>
};


export type QueryFileArgs = {
  id: Scalars['ID']
};


export type QueryFilesArgs = {
  input?: Maybe<Filter>
};


export type QueryCredentialArgs = {
  id: Scalars['ID']
};


export type QueryCredentialsArgs = {
  input?: Maybe<Filter>
};


export type QueryJobArgs = {
  id: Scalars['ID']
};


export type QueryJobsArgs = {
  input?: Maybe<Filter>
};


export type QueryTagArgs = {
  id: Scalars['ID']
};


export type QueryTagsArgs = {
  input?: Maybe<Filter>
};


export type QueryTargetArgs = {
  id: Scalars['ID']
};


export type QueryTargetsArgs = {
  input?: Maybe<Filter>
};


export type QueryTaskArgs = {
  id: Scalars['ID']
};


export type QueryTasksArgs = {
  input?: Maybe<Filter>
};


export type QueryUserArgs = {
  id: Scalars['ID']
};


export type QueryUsersArgs = {
  input?: Maybe<Filter>
};


export type QueryServiceArgs = {
  id: Scalars['ID']
};


export type QueryServicesArgs = {
  input?: Maybe<Filter>
};


export type QueryEventArgs = {
  id: Scalars['ID']
};


export type QueryEventsArgs = {
  input?: Maybe<Filter>
};

export type QueueJobRequest = {
  id: Scalars['ID'],
};

export type RemoveAdminRequest = {
  id: Scalars['ID'],
};

export type RemoveTagRequest = {
  tagID: Scalars['ID'],
  entID: Scalars['ID'],
};

export type Service = {
   __typename?: 'Service',
  id: Scalars['ID'],
  name?: Maybe<Scalars['String']>,
  pubKey?: Maybe<Scalars['String']>,
  config?: Maybe<Scalars['String']>,
  isActivated?: Maybe<Scalars['Boolean']>,
  tag: Tag,
};

export type SetLinkFieldsRequest = {
  id: Scalars['ID'],
  alias?: Maybe<Scalars['String']>,
  ExpirationTime?: Maybe<Scalars['Time']>,
  clicks?: Maybe<Scalars['Int']>,
};

export type SetServiceConfigRequest = {
  id: Scalars['ID'],
  config?: Maybe<Scalars['String']>,
};

export type SetTargetFieldsRequest = {
  id: Scalars['ID'],
  name?: Maybe<Scalars['String']>,
  machineUUID?: Maybe<Scalars['String']>,
  primaryIP?: Maybe<Scalars['String']>,
  publicIP?: Maybe<Scalars['String']>,
  primaryMAC?: Maybe<Scalars['String']>,
  hostname?: Maybe<Scalars['String']>,
};

export type SubmitTaskResultRequest = {
  id: Scalars['ID'],
  output?: Maybe<Scalars['String']>,
  error?: Maybe<Scalars['String']>,
  execStartTime?: Maybe<Scalars['Time']>,
  execStopTime?: Maybe<Scalars['Time']>,
};

export type Tag = {
   __typename?: 'Tag',
  id: Scalars['ID'],
  name?: Maybe<Scalars['String']>,
  tasks?: Maybe<Array<Maybe<Task>>>,
  targets?: Maybe<Array<Maybe<Target>>>,
  jobs?: Maybe<Array<Maybe<Job>>>,
};


export type TagTasksArgs = {
  input?: Maybe<Filter>
};


export type TagTargetsArgs = {
  input?: Maybe<Filter>
};


export type TagJobsArgs = {
  input?: Maybe<Filter>
};

export type Target = {
   __typename?: 'Target',
  id?: Maybe<Scalars['ID']>,
  name?: Maybe<Scalars['String']>,
  primaryIP?: Maybe<Scalars['String']>,
  machineUUID?: Maybe<Scalars['String']>,
  publicIP?: Maybe<Scalars['String']>,
  primaryMAC?: Maybe<Scalars['String']>,
  hostname?: Maybe<Scalars['String']>,
  lastSeen?: Maybe<Scalars['Time']>,
  tasks?: Maybe<Array<Maybe<Task>>>,
  tags?: Maybe<Array<Maybe<Tag>>>,
  credentials?: Maybe<Array<Maybe<Credential>>>,
};


export type TargetTasksArgs = {
  input?: Maybe<Filter>
};


export type TargetTagsArgs = {
  input?: Maybe<Filter>
};


export type TargetCredentialsArgs = {
  input?: Maybe<Filter>
};

export type Task = {
   __typename?: 'Task',
  id: Scalars['ID'],
  lastChangedTime?: Maybe<Scalars['Time']>,
  queueTime?: Maybe<Scalars['Time']>,
  claimTime?: Maybe<Scalars['Time']>,
  execStartTime?: Maybe<Scalars['Time']>,
  execStopTime?: Maybe<Scalars['Time']>,
  content?: Maybe<Scalars['String']>,
  output?: Maybe<Scalars['String']>,
  error?: Maybe<Scalars['String']>,
  sessionID?: Maybe<Scalars['String']>,
  job?: Maybe<Job>,
  target?: Maybe<Target>,
};


export type User = {
   __typename?: 'User',
  id: Scalars['ID'],
  name?: Maybe<Scalars['String']>,
  oAuthID?: Maybe<Scalars['String']>,
  photoURL?: Maybe<Scalars['String']>,
  isActivated?: Maybe<Scalars['Boolean']>,
  isAdmin?: Maybe<Scalars['Boolean']>,
  jobs?: Maybe<Array<Maybe<Job>>>,
};


export type UserJobsArgs = {
  input?: Maybe<Filter>
};


