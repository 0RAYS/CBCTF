import request from './request';

// 获取题目列表
// 请求示例响应:
// {
//    "code": 200,
//    "msg": "Success",
//    "data": [
//        {
//            "challenge": {
//                "category": "Test",
//                "description": "test",
//                "docker": "",
//                "flag": "",
//                "generator_image": "",
//                "id": "21ccf8ad-05cf-497e-896d-ff859fc9ab33",
//                "name": "test",
//                "port": 8080,
//                "type": 0
//            },
//            "usage": {
//                "attempt": 0,
//                "challenge_id": "21ccf8ad-05cf-497e-896d-ff859fc9ab33",
//                "contest_id": 12,
//                "flag": "",
//                "hidden": false,
//                "hints": "",
//                "id": 18,
//                "score": 0,
//                "tags": ""
//            }
//        }
//    ],
//    "trace": "b4a8ce4f-fc6b-4dfb-9457-2c1e76b23171"
// }
export const getChallengeList = (contestId, params = { limit: 20, offset: 0 }) => {
  return request({
    url: `/contests/${contestId}/challenges`,
    method: 'GET',
    params,
  });
};

export const getChallengeCategories = (contestId) => {
  return request({
    url: `/contests/${contestId}/challenges/categories`,
    method: 'GET',
  });
};

// 获取题目状态
// 请求示例响应:
// {
//   "code": 200,
//   "msg": "操作成功",
//   "data": {
//       "files": "",
//       "remote": {
//           "remaining": "",
//           "target": ""
//       },
//       "status": true
//   },
//   "trace": "eae36b06-36bf-4970-a5d3-91399b9f9ed4"
// }
export const getChallengeStatus = (contestId, challengeId) => {
  return request({
    url: `/contests/${contestId}/challenges/${challengeId}`,
    method: 'GET',
  });
};

// 初始化题目
// 请求示例响应:
// {
//    "code": 200,
//    "msg": "Success",
//    "data": null,
//    "trace": "3b1df4fe-befe-4e91-82c8-a702752c7612"
// }
export const initChallenge = (contestId, challengeId) => {
  return request({
    url: `/contests/${contestId}/challenges/${challengeId}/init`,
    method: 'POST',
  });
};

// 下载题目附件
// 响应类型: blob
export const downloadChallengeAttachment = (contestId, challengeId) => {
  return request({
    url: `/contests/${contestId}/challenges/${challengeId}/attachment`,
    method: 'GET',
    responseType: 'blob',
  }).then((response) => {
    return response;
  });
};

// 启动靶机
// 请求示例响应:
// {
//    "code": 200,
//    "msg": "Success",
//    "data": {
//        "remaining": 3598.825783301,
//        "target": "10.0.0.177:31937"
//    },
//    "trace": "02c31e17-73b8-4e7d-9e85-e44d4ad6a134"
// }
export const startRemoteTarget = (contestId, challengeId) => {
  return request({
    url: `/contests/${contestId}/challenges/${challengeId}/start`,
    method: 'POST',
  });
};

// 延长靶机时间
// 请求示例响应:
// {
//    "code": 400,
//    "msg": "Can only extend time within 20 minutes before the container closes",
//    "data": null,
//    "trace": "fb370260-a941-4afa-8953-7a481a861d48"
// }
export const extendContainerTime = (contestId, challengeId) => {
  return request({
    url: `/contests/${contestId}/challenges/${challengeId}/extend`,
    method: 'POST',
  });
};

// 关闭靶机
// 请求示例响应:
// {
//    "code": 200,
//    "msg": "Success",
//    "data": null,
//    "trace": "xxx-xxx-xxx"
// }
export const stopContainer = (contestId, challengeId) => {
  return request({
    url: `/contests/${contestId}/challenges/${challengeId}/stop`,
    method: 'POST',
  });
};

// 提交 flag
// 请求参数示例:
// {
//    "flag": "CBCTF{63983c57-076c-4360-ba9c-1759d426a312}"
// }
// 请求示例响应:
// {
//    "code": 200,
//    "msg": "Success",
//    "data": null,
//    "trace": "d6bdcdd9-db2c-431d-96ba-899b3e508aca"
// }
export const submitFlag = (contestId, challengeId, data) => {
  return request({
    url: `/contests/${contestId}/challenges/${challengeId}/submit`,
    method: 'POST',
    data,
  });
};

export const resetChallenge = (contestId, challengeId) => {
  return request({
    url: `/contests/${contestId}/challenges/${challengeId}/reset`,
    method: 'POST',
  });
};

// 上传题解
// 请求参数为FormData, 包含writeup文件
// 请求示例响应:
// {
//    "code": 200,
//    "msg": "Success",
//    "data": null,
//    "trace": "xxx-xxx-xxx"
// }
export const uploadWriteup = (contestId, file) => {
  const formData = new FormData();
  formData.append('writeup', file);

  return request({
    url: `/contests/${contestId}/writeups`,
    method: 'POST',
    data: formData,
    headers: {
      'Content-Type': 'multipart/form-data',
    },
  });
};

// 获取已上传题解列表
// 请求示例响应:
// {
//    "code": 200,
//    "msg": "Success",
//    "data": {
//      "pictures": [
//        {
//          "date": "2025-04-24T14:14:58.351+08:00",
//          "filename": "stegsolve.docx",
//          "hash": "d8a7fec6d4f31466daea9cec2a3ae2a4b9bf2ff2f6ebdc899c0af63af8f81001",
//          "id": "45c37d57-6a05-418f-b67d-7596d50b8463",
//          "size": 10391,
//          "suffix": ".docx",
//          "uploader": 3
//        }
//      ]
//    },
//    "trace": "xxx-xxx-xxx"
// }
export const getWriteups = (contestId) => {
  return request({
    url: `/contests/${contestId}/writeups`,
    method: 'GET',
  });
};
