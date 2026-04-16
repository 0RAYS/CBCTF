import request from '../request';

// 上传队伍头像
// 请求参数: FormData, 包含picture字段
// 响应示例:
// {
//    "code": 200,
//    "msg": "Success",
//    "data": null,
//    "trace": "xxx-xxx-xxx"
// }
export const uploadTeamPicture = (contestId, file) => {
  const formData = new FormData();
  formData.append('picture', file);
  return request({
    url: `/contests/${contestId}/teams/me/picture`,
    method: 'POST',
    data: formData,
    headers: {
      'Content-Type': 'multipart/form-data',
    },
  });
};
// 解散队伍
export const deleteTeam = (contestId) => {
  return request({
    url: `/contests/${contestId}/teams/me`,
    method: 'DELETE',
  });
};

// 踢出队员
export const kickTeamMember = (contestId, userId) => {
  return request({
    url: `/contests/${contestId}/teams/me/kick`,
    method: 'POST',
    data: {
      user_id: userId,
    },
  });
};

// 获取比赛队伍信息
// 请求示例响应:
// {
//    "code": 200,
//    "msg": "Success",
//    "data": {
//        "picture": "https://test.jbnrz.com.cn/picture/6592c0ad-f64f-4281-a1b4-a5c20f2f13fa",
//        "banned": false,
//        "captain_id": 13,
//        "contest_id": 11,
//        "description": "testsetsetsetset",
//        "hidden": false,
//        "id": 10,
//        "name": "0RAYS",
//        "users": 1
//    },
//    "trace": "ac5f9334-6303-4396-a2bb-32b763fe14ad"
// }
export const getTeamInfo = (contestId) => {
  return request({
    url: `/contests/${contestId}/teams/me`,
    method: 'GET',
  });
};

export const getTeamCaptcha = (contestId) => {
  return request({
    url: `/contests/${contestId}/teams/me/captcha`,
    method: 'GET',
  });
};

export const updateTeamCaptcha = (contestId) => {
  return request({
    url: `/contests/${contestId}/teams/me/captcha`,
    method: 'PUT',
  });
};

// 获取队伍成员列表
// 请求示例响应:
// {
//    "code": 200,
//    "msg": "Success",
//    "data": [
//        {
//            "picture": "https://test.jbnrz.com.cn/picture/ce9c4d0a-9677-4b4e-a8b7-35abcd71702f",
//            "banned": false,
//            "contests": 0,
//            "description": "test",
//            "email": "0rays@0rays.club",
//            "hidden": false,
//            "id": 13,
//            "name": "JBNRZ",
//            "teams": 0,
//            "verified": true
//        }
//    ],
//    "trace": "b7cc140c-b2f6-4c82-b0e2-ad8ace34bed1"
// }
export const getTeamMembers = (contestId) => {
  return request({
    url: `/contests/${contestId}/teams/me/users`,
    method: 'GET',
  });
};

// 更新队伍信息
// 请求参数示例:
// {
//    "name": "test",
//    "description": "testtest",
//    "captcha": "test",
//    "captain_id": 1
// }
export const updateTeamInfo = (contestId, data) => {
  return request({
    url: `/contests/${contestId}/teams/me`,
    method: 'PUT',
    data,
  });
};

// 创建队伍
// 请求参数示例:
// {
//    "name": "test-team",
//    "description": "testsetsetset",
//    "captcha": ""
// }
export const createTeam = (contestId, data) => {
  return request({
    url: `/contests/${contestId}/teams/create`,
    method: 'POST',
    data,
  });
};

// 加入队伍
// 请求参数示例:
// {
//    "name": "0RAYS",
//    "captcha": "fake-captcha"
// }
export const joinTeam = (contestId, data) => {
  return request({
    url: `/contests/${contestId}/teams/join`,
    method: 'POST',
    data,
  });
};
