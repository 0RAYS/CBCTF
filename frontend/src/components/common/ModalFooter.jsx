import Button from './Button';

function ModalFooter({
  onCancel,
  onSubmit,
  cancelLabel,
  submitLabel,
  submitVariant = 'primary',
  submitDisabled = false,
  submitLoading = false,
}) {
  return (
    <>
      <Button size="sm" variant="ghost" onClick={onCancel} disabled={submitLoading}>
        {cancelLabel}
      </Button>
      <Button size="sm" variant={submitVariant} onClick={onSubmit} disabled={submitDisabled} loading={submitLoading}>
        {submitLabel}
      </Button>
    </>
  );
}

export default ModalFooter;
