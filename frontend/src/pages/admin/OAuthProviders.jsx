import { useEffect, useState } from 'react';
import { useTranslation } from 'react-i18next';
import { getOAuthProviderList } from '../../api/admin/oauth';
import AdminOAuthProviders from '../../components/features/Admin/oauth/AdminOAuthProviders';
import OAuthProviderDialog from '../../components/features/Admin/oauth/OAuthProviderDialog';
import useProviderPicture from '../../components/features/Admin/oauth/useProviderPicture';
import { toast } from '../../utils/toast';

export default function OAuthProvidersManagement() {
  const { t } = useTranslation();
  const [providers, setProviders] = useState([]);
  const [totalCount, setTotalCount] = useState(0);
  const [currentPage, setCurrentPage] = useState(1);
  const [dialog, setDialog] = useState(null);
  const pageSize = 20;

  const fetchProviders = async () => {
    try {
      const response = await getOAuthProviderList({ limit: pageSize, offset: (currentPage - 1) * pageSize });
      if (response.code === 200) {
        setProviders(response.data.providers);
        setTotalCount(response.data.count);
      }
    } catch (error) {
      toast.danger({ description: error.message || t('admin.oauthProviders.toast.fetchFailed') });
    }
  };

  useEffect(() => {
    fetchProviders();
  }, [currentPage]);
  const { inputRef, open: uploadPicture, onChange: handlePictureChange } = useProviderPicture(fetchProviders);
  const edit = (provider) => setDialog({ mode: 'edit', provider });

  return (
    <>
      <AdminOAuthProviders
        providers={providers}
        totalCount={totalCount}
        currentPage={currentPage}
        pageSize={pageSize}
        onPageChange={setCurrentPage}
        onCreateProvider={() => setDialog({ mode: 'create' })}
        onEditProvider={edit}
        onDeleteProvider={(provider) => setDialog({ mode: 'delete', provider })}
        onProviderClick={edit}
        onPictureUpload={uploadPicture}
      />
      {dialog && (
        <OAuthProviderDialog
          {...dialog}
          onClose={() => setDialog(null)}
          onSaved={() => {
            setDialog(null);
            fetchProviders();
          }}
        />
      )}
      <input
        type="file"
        ref={inputRef}
        className="hidden"
        accept="image/png,image/jpeg,image/jpg,image/gif"
        onChange={handlePictureChange}
      />
    </>
  );
}
