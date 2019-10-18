import {Url} from '../entity/Url';
import {ApolloClient} from 'apollo-client';
import {HttpLink} from 'apollo-link-http';
import {InMemoryCache, NormalizedCacheObject} from 'apollo-cache-inmemory';
import {ApolloLink, FetchResult} from 'apollo-link';
import gql from 'graphql-tag';
import {EnvService} from './Env.service';
import {GraphQlError} from '../graphql/error';
import {AuthService} from './Auth.service';
import {CaptchaService, CREATE_SHORT_LINK} from './Captcha.service';
import {validateLongLinkFormat} from '../validators/LongLink.validator';
import {validateCustomAliasFormat} from '../validators/CustomAlias.validator';
import {ErrorService, ErrUrl} from './Error.service';

interface CreateURLData {
  createURL: Url;
}

export class UrlService {
  private gqlClient: ApolloClient<NormalizedCacheObject>;

  constructor(
    private authService: AuthService,
    private envService: EnvService,
    private errorService: ErrorService,
    private captchaService: CaptchaService
  ) {
    const gqlLink = ApolloLink.from([
      new HttpLink({
        uri: `${this.envService.getVal('GRAPHQL_API_BASE_URL')}/graphql`
      })
    ]);

    this.gqlClient = new ApolloClient({
      link: gqlLink,
      cache: new InMemoryCache()
    });
  }

  createShortLink(editingUrl: Url): Promise<Url> {
    return new Promise(async (resolve, reject) => {
      let longLink = editingUrl.originalUrl;
      let customAlias = editingUrl.alias;

      let err = validateLongLinkFormat(longLink);
      if (err) {
        reject({
          createShortLinkErr: {
            name: 'Invalid Long Link',
            description: err
          }
        });
        return;
      }

      err = validateCustomAliasFormat(customAlias);
      if (err) {
        reject({
          createShortLinkErr: {
            name: 'Invalid Custom Alias',
            description: err
          }
        });
        return;
      }

      try {
        let url = await this.invokeCreateShortLinkApi(
          editingUrl
        );
        resolve(url);
        return;
      } catch (errCodes) {
        let errCode = errCodes[0];
        if (errCode === ErrUrl.Unauthorized) {
          reject({
            authorizationErr: 'Unauthorized to create short link'
          });
          return;
        }

        let error = this.errorService.getErr(errCode);
        reject({
          createShortLinkErr: error
        });
      }
    });
  }

  private async invokeCreateShortLinkApi(link: Url): Promise<Url> {
    // let captchaResponse = await this.captchaService.execute(CREATE_SHORT_LINK);

    let alias = link.alias === '' ? null : link.alias;

    let variables = {
      captchaResponse: 'null',
      urlInput: {
        originalURL: link.originalUrl,
        customAlias: alias
      },
      authToken: this.authService.getAuthToken()
    };

    let mutation = gql`
      mutation params(
        $captchaResponse: String!
        $urlInput: URLInput!
        $authToken: String!
      ) {
        createURL(
          captchaResponse: $captchaResponse
          url: $urlInput
          authToken: $authToken
        ) {
          alias
          originalURL
        }
      }
    `;

    return new Promise<Url>((resolve, reject: (errCodes: ErrUrl[]) => any) => {
      this.gqlClient
        .mutate({
          variables: variables,
          mutation: mutation
        })
        .then((res: FetchResult<CreateURLData>) => {
          if (!res || !res.data) {
            return resolve({});
          }
          resolve(res.data.createURL);
        })
        .catch(({graphQLErrors, networkError, message}) => {
          const errCodes = graphQLErrors.map(
            (graphQLError: GraphQlError) => graphQLError.extensions.code
          );
          reject(errCodes);
        });
    });
  }

  aliasToLink(alias: string): string {
    return `${this.envService.getVal('HTTP_API_BASE_URL')}/r/${alias}`;
  }
}
