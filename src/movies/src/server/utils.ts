import {httpResponse} from "../model/model";

export function createSuccessResponse(data: any): httpResponse {
    return {
        status: true,
        data: data
    }
}

export function createErrorResponse(data: any): httpResponse {
    return {
        status: false,
        data: data
    }
}