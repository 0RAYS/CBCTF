import { useState, useEffect, useCallback, useMemo, useRef } from 'react';
import { useTranslation } from 'react-i18next';
import { IconSearch, IconSortAscending, IconSortDescending } from '@tabler/icons-react';
import { getSearchModels, searchModels } from '../../../api/admin/search';
import { useDebounce } from '../../../hooks/index.js';
import { toast } from '../../../utils/toast';
import Modal from '../../common/Modal';
import Select from '../../common/Select';
import Input from '../../common/Input';
import Button from '../../common/Button';
import Loading from '../../common/Loading';
import EmptyState from '../../common/EmptyState';
import Pagination from '../../common/Pagination';

const PAGE_SIZE = 20;
const CELL_MAX_LENGTH = 50;

function formatValue(value) {
  if (value === null || value === undefined) return '-';
  if (typeof value === 'boolean') return value ? 'Yes' : 'No';
  if (typeof value === 'object') return JSON.stringify(value, null, 2);
  return String(value);
}

function cellText(value) {
  if (value === null || value === undefined) return '-';
  if (typeof value === 'boolean') return value ? 'Yes' : 'No';
  if (typeof value === 'object') {
    const s = JSON.stringify(value);
    return s.length > CELL_MAX_LENGTH ? s.slice(0, CELL_MAX_LENGTH) + '...' : s;
  }
  const s = String(value);
  return s.length > CELL_MAX_LENGTH ? s.slice(0, CELL_MAX_LENGTH) + '...' : s;
}

function AdminGlobalSearch({ isOpen, onClose }) {
  const { t } = useTranslation();

  const [modelsMap, setModelsMap] = useState({});
  const [selectedModel, setSelectedModel] = useState('');
  const [detailItem, setDetailItem] = useState(null);
  const [filters, setFilters] = useState({});
  const [sorts, setSorts] = useState({});
  const [results, setResults] = useState([]);
  const [total, setTotal] = useState(0);
  const [page, setPage] = useState(1);
  const [loading, setLoading] = useState(false);

  const debouncedFilters = useDebounce(filters, 400);
  const abortRef = useRef(null);

  useEffect(() => {
    if (!isOpen) return;
    getSearchModels()
      .then((res) => {
        if (res.code === 200 && res.data) {
          setModelsMap(res.data);
        } else {
          toast.warning({ description: t('admin.globalSearch.toast.modelsFailed') });
        }
      })
      .catch(() => {
        toast.warning({ description: t('admin.globalSearch.toast.modelsFailed') });
      });
  }, [isOpen, t]);

  useEffect(() => {
    setFilters({});
    setSorts({});
    setResults([]);
    setTotal(0);
    setPage(1);
  }, [selectedModel]);

  const executeSearch = useCallback(async () => {
    if (!selectedModel) return;

    if (abortRef.current) abortRef.current.abort();
    abortRef.current = new AbortController();

    const params = {
      model: selectedModel,
      limit: PAGE_SIZE,
      offset: (page - 1) * PAGE_SIZE,
    };

    for (const [field, value] of Object.entries(debouncedFilters)) {
      if (value) {
        params[`search[${field}]`] = value;
      }
    }

    for (const [field, dir] of Object.entries(sorts)) {
      if (dir) {
        params[`sort[${field}]`] = dir;
      }
    }

    setLoading(true);
    try {
      const res = await searchModels(params);
      if (res.code === 200 && res.data) {
        setResults(res.data.models || []);
        setTotal(res.data.count || 0);
      } else {
        toast.warning({ description: t('admin.globalSearch.toast.searchFailed') });
      }
    } catch (err) {
      if (err?.name !== 'CanceledError' && err?.name !== 'AbortError') {
        toast.warning({ description: t('admin.globalSearch.toast.searchFailed') });
      }
    } finally {
      setLoading(false);
    }
  }, [selectedModel, debouncedFilters, sorts, page, t]);

  useEffect(() => {
    if (selectedModel) {
      executeSearch();
    }
  }, [executeSearch, selectedModel]);

  useEffect(() => {
    setPage(1);
  }, [debouncedFilters, sorts]);

  useEffect(() => {
    if (!isOpen) {
      setSelectedModel('');
      setFilters({});
      setSorts({});
      setResults([]);
      setTotal(0);
      setPage(1);
    }
  }, [isOpen]);

  const selectedModelConfig = modelsMap[selectedModel] || {};
  const queryFields = Array.isArray(selectedModelConfig) ? selectedModelConfig : selectedModelConfig.query || [];
  const searchFields = Array.isArray(selectedModelConfig) ? selectedModelConfig : selectedModelConfig.search || [];

  const modelOptions = Object.keys(modelsMap).map((name) => ({
    value: name,
    label: name,
  }));

  const handleFilterChange = (field, value) => {
    setFilters((prev) => ({ ...prev, [field]: value }));
  };

  const handleSortToggle = (field) => {
    setSorts((prev) => {
      const current = prev[field];
      let next;
      if (!current) next = 'asc';
      else if (current === 'asc') next = 'desc';
      else next = null;

      const updated = { ...prev };
      if (next) {
        updated[field] = next;
      } else {
        delete updated[field];
      }
      return updated;
    });
  };

  const getSortIcon = (field) => {
    const dir = sorts[field];
    if (dir === 'asc') return <IconSortAscending size={16} />;
    if (dir === 'desc') return <IconSortDescending size={16} />;
    return null;
  };

  const columnKeys = useMemo(() => {
    if (results.length === 0) return [];
    return Object.keys(results[0]);
  }, [results]);

  const totalPages = Math.max(1, Math.ceil(total / PAGE_SIZE));

  return (
    <>
      {/* Search modal — use portal directly to bypass Modal max-width */}
      <Modal
        isOpen={isOpen}
        onClose={onClose}
        title={t('admin.globalSearch.title')}
        size="2xl"
        className="max-w-[95vw]!"
      >
        <div className="space-y-4">
          {/* Model selector */}
          <div>
            <label className="block text-sm text-neutral-400 font-mono mb-1">{t('admin.globalSearch.model')}</label>
            <Select
              value={selectedModel}
              onChange={(e) => setSelectedModel(e.target.value)}
              options={modelOptions}
              placeholder={t('admin.globalSearch.selectModel')}
              size="sm"
            />
          </div>

          {/* Dynamic filter fields */}
          {searchFields.length > 0 && (
            <div>
              <label className="block text-sm text-neutral-400 font-mono mb-1">{t('admin.globalSearch.filters')}</label>
              <div className="grid grid-cols-2 md:grid-cols-3 gap-2">
                {searchFields.map((field) => (
                  <Input
                    key={field}
                    size="sm"
                    placeholder={field}
                    value={filters[field] || ''}
                    onChange={(e) => handleFilterChange(field, e.target.value)}
                    icon={<IconSearch size={14} />}
                  />
                ))}
              </div>
            </div>
          )}

          {/* Sort controls */}
          {queryFields.length > 0 && (
            <div>
              <label className="block text-sm text-neutral-400 font-mono mb-1">{t('admin.globalSearch.sort')}</label>
              <div className="flex flex-wrap gap-1">
                {queryFields.map((field) => (
                  <Button
                    key={field}
                    size="sm"
                    variant={sorts[field] ? 'primary' : 'ghost'}
                    onClick={() => handleSortToggle(field)}
                    icon={getSortIcon(field)}
                  >
                    {field}
                  </Button>
                ))}
              </div>
            </div>
          )}

          {/* Results table */}
          {selectedModel && (
            <div>
              {loading ? (
                <Loading />
              ) : results.length === 0 ? (
                <EmptyState title={t('admin.globalSearch.noResults')} />
              ) : (
                <div className="overflow-x-auto border border-neutral-300/10 rounded">
                  <table className="w-full" style={{ tableLayout: 'auto' }}>
                    <thead>
                      <tr className="bg-black/40">
                        {columnKeys.map((key) => (
                          <th
                            key={key}
                            className="px-3 py-2 text-left text-xs text-neutral-400 font-mono whitespace-nowrap"
                          >
                            {key}
                          </th>
                        ))}
                      </tr>
                    </thead>
                    <tbody>
                      {results.map((item, rowIndex) => (
                        <tr
                          key={rowIndex}
                          className="border-t border-neutral-300/10 hover:bg-black/40 transition-colors cursor-pointer"
                          onClick={() => setDetailItem(item)}
                        >
                          {columnKeys.map((key) => (
                            <td
                              key={key}
                              className="px-3 py-2 text-sm text-neutral-300 font-mono whitespace-nowrap max-w-65 overflow-hidden text-ellipsis"
                              title={formatValue(item[key])}
                            >
                              {cellText(item[key])}
                            </td>
                          ))}
                        </tr>
                      ))}
                    </tbody>
                  </table>
                </div>
              )}

              {total > PAGE_SIZE && (
                <div className="pt-3">
                  <Pagination current={page} total={totalPages} onChange={setPage} showTotal totalItems={total} />
                </div>
              )}
            </div>
          )}
        </div>
      </Modal>

      {/* Row detail modal */}
      <Modal
        isOpen={!!detailItem}
        onClose={() => setDetailItem(null)}
        title={t('admin.globalSearch.detailTitle', { model: selectedModel })}
        size="lg"
      >
        {detailItem && (
          <div className="space-y-3">
            {Object.entries(detailItem).map(([key, value]) => (
              <div key={key} className="border-b border-neutral-300/10 pb-3 last:border-b-0">
                <div className="text-xs text-neutral-500 font-mono mb-1">{key}</div>
                <pre className="text-sm text-neutral-200 font-mono whitespace-pre-wrap break-all">
                  {formatValue(value)}
                </pre>
              </div>
            ))}
          </div>
        )}
      </Modal>
    </>
  );
}

export default AdminGlobalSearch;
