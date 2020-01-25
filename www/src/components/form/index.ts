import { useState } from "react";

export { default as XScriptEditor } from "./XScriptEditor";
export { default as XTagTypeahead, SUGGEST_TAGS_QUERY } from "./XTagTypeahead";
export {
  default as XTargetTypeahead,
  SUGGEST_TARGETS_QUERY
} from "./XTargetTypeahead";

export const useModal = (): [() => void, () => void, boolean] => {
  const [isOpen, setOpen] = useState(false);

  const openModal = () => setOpen(true);
  const closeModal = () => setOpen(false);

  return [openModal, closeModal, isOpen];
};
