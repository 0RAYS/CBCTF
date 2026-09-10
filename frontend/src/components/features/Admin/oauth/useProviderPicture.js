import { useRef } from 'react';
import { useTranslation } from 'react-i18next';
import { uploadOAuthPicture } from '../../../../api/admin/oauth';
import { toast } from '../../../../utils/toast';

export default function useProviderPicture(onUploaded) {
  const { t } = useTranslation();
  const inputRef = useRef(null);
  const targetRef = useRef(null);

  const open = (provider) => {
    targetRef.current = provider;
    inputRef.current?.click();
  };
  const onChange = async (event) => {
    const file = event.target.files?.[0];
    const target = targetRef.current;
    event.target.value = '';
    if (!file || !target) return;
    try {
      const response = await uploadOAuthPicture(target.id, file);
      if (response.code === 200) {
        toast.success({ description: t('admin.oauthProviders.toast.pictureUpdated') });
        onUploaded();
      }
    } catch (error) {
      toast.danger({ description: error.message || t('admin.oauthProviders.toast.pictureUpdateFailed') });
    }
  };
  return { inputRef, open, onChange };
}
