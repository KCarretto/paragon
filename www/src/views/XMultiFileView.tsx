import { useQuery } from "@apollo/react-hooks";
import gql from "graphql-tag";
import * as React from "react";
import { XFileCard, XNoFilesFound } from "../components/file";
import { XBoundary, XCardGroup } from "../components/layout";
import { XErrorMessage, XLoadingMessage } from "../components/messages";
import { File } from "../graphql/models";

export const MULTI_FILE_QUERY = gql`
  {
    files {
      id
      name
      contentType
      size
      creationTime
      lastModifiedTime

      links {
        id
        alias
        clicks
        expirationTime
      }
    }
  }
`;

type MultiFile = {
  files: File[];
};

const XMultiFileView = () => {
  const { loading, error, data: { files = [] } = {} } = useQuery<MultiFile>(
    MULTI_FILE_QUERY,
    {
      pollInterval: 5000
    }
  );

  const whenLoading = (
    <XLoadingMessage title="Loading Files" msg="Fetching file list" />
  );
  const whenEmpty = <XNoFilesFound />;

  return (
    <React.Fragment>
      <XErrorMessage title="Error Loading Files" err={error} />
      <XBoundary boundary={whenLoading} show={!loading}>
        <XBoundary boundary={whenEmpty} show={files && files.length > 0}>
          <XCardGroup>
            {files && files.map(file => <XFileCard key={file.id} {...file} />)}
          </XCardGroup>
        </XBoundary>
      </XBoundary>
    </React.Fragment>
  );
};

export default XMultiFileView;
