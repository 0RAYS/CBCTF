import Notices from '../../components/features/CTFGame/Notice/Notice';
import { useState, useEffect } from 'react';
import { getContestNotices } from '../../api/contest';
import { useParams } from 'react-router-dom';

const transformTimestamp = (timestamp) => {
  const date = new Date(timestamp);
  return date.toLocaleString();
};

function GameNoticePage() {
  const { contestId } = useParams();
  const [notices, setNotices] = useState([]);

  useEffect(() => {
    getContestNotices(contestId).then((res) => {
      setNotices(
        res.data.notices.map((notice) => ({
          id: notice.id,
          title: notice.title,
          content: notice.content,
          type: notice.type || 'normal',
          timestamp: transformTimestamp(notice.created_at),
        }))
      );
    });
  }, [contestId]);

  return <Notices notices={notices} />;
}

export default GameNoticePage;
