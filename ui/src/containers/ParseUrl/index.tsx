import React from 'react';
import { useNavigate } from 'react-router-dom';
import { useMst } from '../../store/root';
import { Params } from '../../common/params';
import { AuthCodeProps, IError } from '../../store/auth';

const ParseUrl: React.FC = () => {
  const { resources, user } = useMst();
  const history = useNavigate();

  if (window.location.search) {
    const searchParams: URLSearchParams = new URLSearchParams(window.location.search);
    const status = searchParams.get('status');
    const code = searchParams.get('code');

    // It checks status and code and then redirect to authentication
    if (status === '200' && code !== null) {
      const codeFinal: AuthCodeProps = {
        code: code
      };
      user.authenticate(codeFinal);
      if (user.isAuthenticated) {
        // Initially `history.goBack` was used to go to the previous
        // page but with the update of `react-router-dom` version '-1'
        // is added so that the page redirects back to the previous page
        history('-1');
      }
    }
    // Display the alert message when status is not ok
    else if (!user.isAuthenticated && status !== '200' && status !== null) {
      // Wait to redirection of page and then update the store
      setTimeout(() => {
        const error: IError = {
          status: Number(status),
          serverMessage: 'Login Failed, Please Try To Login Again!',
          customMessage: ''
        };
        user.setErrorMessage(error);
      }, 1000);
    }
    if (searchParams.has(Params.Query)) {
      resources.setSearch(searchParams.get(Params.Query) || '');
    }
    if (searchParams.has(Params.Tag)) {
      const tags = searchParams.getAll(Params.Tag);
      resources.setSearch(`tags:${tags.join(',')}`);
      resources.setSearchedTags(searchParams.getAll(Params.Tag));
    }
    if (searchParams.has(Params.SortBy)) {
      resources.setSortBy(searchParams.get(Params.SortBy) || '');
    }
    // Storing url params to store inorder to parse the url only after successfully resource load
    resources.setURLParams(window.location.search);
  }
  return <> </>;
};
export default ParseUrl;
