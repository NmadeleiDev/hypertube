import { httpResponse, IMovie } from '../model/model';

export function createSuccessResponse(data: any): httpResponse {
  return {
    status: true,
    data: data,
  };
}

export function createErrorResponse(data: any): httpResponse {
  return {
    status: false,
    data: data,
  };
}

export const getUserIdFromToken = (token: string) => {
  try {
    return +JSON.parse(
      Buffer.from(
        Buffer.from(token, 'base64').toString().split('.')[0],
        'base64'
      ).toString()
    ).userId;
  } catch (e) {
    return 0;
  }
};

export const isIMovie = (movie: IMovie | unknown): movie is IMovie => {
  return (
    (movie as IMovie).id !== undefined && (movie as IMovie).info !== undefined
  );
};
