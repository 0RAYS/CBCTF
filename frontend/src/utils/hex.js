import i18n from '../i18n';

// 将hex字符串转换为UTF-8字符串
export function hexToUtf8(hex) {
  try {
    // 移除可能存在的空格
    hex = hex.replace(/\s/g, '');

    // 检查是否是有效的hex字符串
    if (!/^[0-9A-Fa-f]*$/.test(hex)) {
      return i18n.t('utils.hex.invalid');
    }

    // 将hex转换为字节数组
    const bytes = new Uint8Array(hex.match(/.{1,2}/g).map((byte) => parseInt(byte, 16)));

    // 使用TextDecoder将字节数组转换为UTF-8字符串
    const decoder = new TextDecoder('utf-8');
    return decoder.decode(bytes);
  } catch (error) {
    return error.message;
  }
}
