import ModalButton from './ModalButton';

function ModalFooter({
  onCancel,
  onSubmit,
  cancelLabel,
  submitLabel,
  submitVariant = 'primary',
  submitDisabled = false,
}) {
  return (
    <>
      <ModalButton onClick={onCancel}>{cancelLabel}</ModalButton>
      <ModalButton variant={submitVariant} onClick={onSubmit} disabled={submitDisabled}>
        {submitLabel}
      </ModalButton>
    </>
  );
}

export default ModalFooter;
