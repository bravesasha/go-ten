import { ResponseDataInterface } from "@/src/types";
import { httpRequest } from "@/api";
import { pathToUrl } from "@/src/routes/router";
import { apiRoutes } from "@/src/routes";

export const fetchTestnetStatus = async (): Promise<
  ResponseDataInterface<boolean>
> => {
  return await httpRequest<ResponseDataInterface<boolean>>({
    method: "get",
    url: pathToUrl(apiRoutes.getHealthStatus),
  });
};