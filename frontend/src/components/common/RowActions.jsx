import Button from './Button';
import TruncatedText from './TruncatedText';

/** The column and context menu share action callbacks, visibility and disabled state. */
export default function RowActions({ actions }) {
  return (
    <div className="flex flex-wrap items-center gap-2" onClick={(event) => event.stopPropagation()}>
      {actions
        .filter((action) => action && action.inline && !action.hidden)
        .map((action) => (
          <Button
            key={action.key}
            variant={action.danger ? 'danger' : 'ghost'}
            size={action.icon ? 'icon' : 'sm'}
            className="min-w-0 max-w-full shrink-0"
            title={action.label}
            aria-label={action.label}
            disabled={action.disabled}
            onClick={() => action.onClick()}
          >
            {action.icon || <TruncatedText>{action.label}</TruncatedText>}
          </Button>
        ))}
    </div>
  );
}
