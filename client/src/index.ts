import axios, { AxiosInstance } from "axios";

var instance: AxiosInstance;

export function init(baseURL: string) {
  instance = axios.create({
    baseURL,
  });
}

/**
 *  根据id 照抄数据
 * collection: shop/category;
 * id: 1111;
 */
export async function findById({
  collection,
  id,
}: {
  collection: string;
  id: string;
}) {
  // const res = await httpRequest.get(`/api/collection/${collection}/${id}`);
  const res = await instance.get(`/api/collection/${collection}/${id}`);
  return res.data;
}

export async function queryForPage({
  collection,
  aggregate,
  curPage,
  pageSize,
}: {
  collection: string;
  aggregate: any[];
  curPage: number;
  pageSize: number;
}) {
  return await instance.post(`/api/collection/query/${collection}`, {
    page: {
      curPage,
      pageSize,
    },
    aggregate: JSON.stringify(aggregate),
  });
}

export async function findOne({
  collection,
  filter,
}: {
  collection: string;
  filter: any;
}) {
  return await instance.post(`/api/collection/findone/${collection}`, {
    filter: filter || {},
    sort: [
      {
        key: "sort",
        sort: -1,
      },
    ],
  });
}

/**
 * 保存/更新
 * @param {*} param0
 * @returns
 */
export async function save({
  body,
  collection,
  filter,
}: {
  body: any;
  collection: string;
  filter: any;
}) {
  const config: any = { body: body };
  if (filter) {
    config.filter = filter;
  }
  return await instance.post(`/api/collection/save/${collection}`, config);
}

export async function delById({
  id,
  collection,
}: {
  id: string;
  collection: string;
}) {
  return await instance.post(`/api/collection/delete/${collection}/${id}`);
}

export async function func({
  funcName,
  collection,
  method,
  config,
}: {
  funcName: string;
  collection: string;
  method: string;
  config: any;
}) {
  if (method.toLowerCase() == "get") {
    return await instance.get(
      `/api/collection/func/${collection}/${funcName}`,
      config || {}
    );
  } else {
    return await instance.post(
      `/api/collection/func/${collection}/${funcName}`,
      config || {}
    );
  }
}
