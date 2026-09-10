import { useEffect, useState } from 'react';
import { useTranslation } from 'react-i18next';
import { getContestChallenges } from '../../../../api/admin/contest';
import { toast } from '../../../../utils/toast';
import { updateChallengeSelection } from './victimPayload';

export default function useVictimCandidates(contestId) {
  const { t } = useTranslation();
  const [challenges, setChallenges] = useState([]);
  const [challengeTotal, setChallengeTotal] = useState(0);
  const [challengePage, setChallengePage] = useState(1);
  const [challengeSearch, setChallengeSearch] = useState('');
  const [detailChallenges, setDetailChallenges] = useState([]);
  const [detailChallengeTotal, setDetailChallengeTotal] = useState(0);
  const [detailChallengePage, setDetailChallengePage] = useState(1);
  const [detailsOpen, setDetailsOpen] = useState(false);
  const [selectedChallenges, setSelectedChallenges] = useState([]);
  const [knownChallenges, setKnownChallenges] = useState({});
  const challengePageSize = 30;

  useEffect(() => {
    let active = true;
    const fetchPage = async (page, setItems, setTotal) => {
      try {
        const response = await getContestChallenges(contestId, {
          type: 'pods',
          limit: challengePageSize,
          offset: (page - 1) * challengePageSize,
          ...(challengeSearch.trim() && { name: challengeSearch.trim() }),
        });
        if (!active || response.code !== 200) return;
        const items = response.data.challenges || [];
        setItems(items);
        setTotal(response.data.count || 0);
        // Keep labels for selections made on other pages or in the expanded picker.
        setKnownChallenges((previous) => ({
          ...previous,
          ...Object.fromEntries(items.map((item) => [item.id, item])),
        }));
      } catch (error) {
        if (active)
          toast.danger({ description: error.message || t('admin.contests.containers.toast.fetchChallengesFailed') });
      }
    };
    fetchPage(challengePage, setChallenges, setChallengeTotal);
    if (detailsOpen) fetchPage(detailChallengePage, setDetailChallenges, setDetailChallengeTotal);
    return () => {
      active = false;
    };
  }, [contestId, challengePage, challengeSearch, detailsOpen, detailChallengePage, t]);

  return {
    challenges,
    challengeTotal,
    challengePage,
    challengePageSize,
    challengeSearch,
    detailChallenges,
    detailChallengeTotal,
    detailChallengePage,
    detailsOpen,
    selectedChallenges,
    setSelectedChallenges,
    selectedChallengeDetails: selectedChallenges.map((id) => knownChallenges[id]).filter(Boolean),
    setChallengePage,
    setDetailChallengePage,
    openDetails: () => {
      setDetailChallengePage(1);
      setDetailsOpen(true);
    },
    closeDetails: () => setDetailsOpen(false),
    changeSearch: (value) => {
      setChallengeSearch(value);
      setChallengePage(1);
      setDetailChallengePage(1);
    },
    selectChallenge: (id, checked) =>
      setSelectedChallenges((previous) => updateChallengeSelection(previous, id, checked)),
  };
}
