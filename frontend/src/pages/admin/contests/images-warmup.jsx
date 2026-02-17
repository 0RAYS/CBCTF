import { useState, useEffect } from 'react';
import { useParams } from 'react-router-dom';
import { toast } from '../../../utils/toast';
import { getContestWarmupImages, warmupContestImages } from '../../../api/admin/contest';
import AdminImagesWarmup from '../../../components/features/Admin/Contests/AdminImagesWarmup.jsx';
import { useTranslation } from 'react-i18next';

function AdminContestImagesWarmup() {
  const { id } = useParams();
  const [loading, setLoading] = useState(true);
  const [submitting, setSubmitting] = useState(false);
  const [images, setImages] = useState([]);
  const [selectedImages, setSelectedImages] = useState([]);
  const [pullPolicy, setPullPolicy] = useState('IfNotPresent');
  const { t } = useTranslation();

  useEffect(() => {
    fetchWarmupImages();
  }, [id]);

  const fetchWarmupImages = async () => {
    setLoading(true);
    try {
      const response = await getContestWarmupImages(parseInt(id));
      if (response.code === 200) {
        setImages(response.data || []);
      }
    } catch (error) {
      toast.danger({ description: error.message || t('admin.contests.imagesWarmup.toast.fetchFailed') });
    } finally {
      setLoading(false);
    }
  };

  const handleImageToggle = (imageName) => {
    setSelectedImages((prev) =>
      prev.includes(imageName) ? prev.filter((img) => img !== imageName) : [...prev, imageName]
    );
  };

  const handleSelectAll = () => {
    if (selectedImages.length === images.length) {
      setSelectedImages([]);
    } else {
      setSelectedImages(images.map((imageObj) => Object.keys(imageObj)[0]));
    }
  };

  const handleWarmup = async () => {
    if (selectedImages.length === 0) {
      toast.warning({ description: t('admin.contests.imagesWarmup.toast.selectRequired') });
      return;
    }

    setSubmitting(true);
    try {
      const response = await warmupContestImages(parseInt(id), {
        images: selectedImages,
        pull_policy: pullPolicy,
      });

      if (response.code === 200) {
        toast.success({ description: t('admin.contests.imagesWarmup.toast.submitSuccess') });
        // 重新获取数据以更新状态
        setTimeout(() => {
          fetchWarmupImages();
        }, 2000);
      }
    } catch (error) {
      toast.danger({ description: error.message || t('admin.contests.imagesWarmup.toast.warmupFailed') });
    } finally {
      setSubmitting(false);
    }
  };

  return (
    <AdminImagesWarmup
      images={images}
      selectedImages={selectedImages}
      pullPolicy={pullPolicy}
      loading={loading}
      submitting={submitting}
      onImageToggle={handleImageToggle}
      onSelectAll={handleSelectAll}
      onPullPolicyChange={setPullPolicy}
      onWarmup={handleWarmup}
      onRefresh={fetchWarmupImages}
    />
  );
}

export default AdminContestImagesWarmup;
