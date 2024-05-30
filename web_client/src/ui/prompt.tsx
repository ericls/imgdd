import React from "react";
import { Fragment } from "react";
import { Dialog, Transition } from "@headlessui/react";
import { Button } from "./button";
import { noop } from "lodash-es";

type PromptProps = {
  title: string;
  content: React.ReactElement;
  yesText?: string;
  yesDestructive?: boolean;
  onClose: () => void;
  onYes: () => void;
  onNo: () => void;
};
type PromptStackItem = Omit<PromptProps, "onClose" | "onNo" | "onYes"> & {
  id: string;
  onResolve: (value: boolean | null) => void;
};

class PromptStack extends EventTarget {
  prompts: PromptStackItem[] = [];
  push = (newPrompt: PromptStackItem) => {
    this.prompts.push(newPrompt);
    this.dispatchEvent(new Event("changed"));
  };
  remove = (id: string) => {
    const existingIndex = this.prompts.findIndex((p) => p.id === id);
    if (existingIndex > -1) {
      this.prompts.splice(existingIndex, 1);
    }
    this.dispatchEvent(new Event("changed"));
  };
}

const _promptStack = new PromptStack();

export function PromptContainer() {
  const [_flag, setFlag] = React.useState(false);
  React.useEffect(() => {
    const cb = () => {
      setFlag((v) => !v);
    };
    _promptStack.addEventListener("changed", cb);
    return () => {
      _promptStack.removeEventListener("changed", cb);
    };
  }, [setFlag]);
  const children = _promptStack.prompts.map((pItem) => {
    const { content, yesDestructive, yesText, title, onResolve, id } = pItem;
    const makeCloseFunction = (func: () => void) => {
      return () => {
        func();
        _promptStack.remove(id);
      };
    };
    const onYes = makeCloseFunction(() => onResolve(true));
    const onNo = makeCloseFunction(() => onResolve(false));
    const onClose = makeCloseFunction(() => onResolve(null));
    return (
      <Prompt
        key={id}
        content={content}
        onClose={onClose}
        onYes={onYes}
        onNo={onNo}
        title={title}
        yesText={yesText}
        yesDestructive={yesDestructive}
      />
    );
  });
  return <>{children}</>;
}

function Prompt({
  title,
  content,
  yesText = "OK",
  yesDestructive = false,
  onClose: _onClose,
  onYes,
  onNo,
}: PromptProps) {
  return (
    // <Transition
    //   show
    //   enter="transition duration-100 ease-out"
    //   enterFrom="transform scale-95 opacity-0"
    //   enterTo="transform scale-100 opacity-100"
    //   leave="transition duration-75 ease-out"
    //   leaveFrom="transform scale-100 opacity-100"
    //   leaveTo="transform scale-95 opacity-0"
    //   as={Fragment}
    // >
    //   <Dialog onClose={noop} as="div" className="relative z-10">
    //     <Dialog.Panel className="relative transform overflow-hidden rounded-lg bg-white text-left shadow-xl transition-all sm:my-8 sm:w-full sm:max-w-lg">
    //       <Dialog.Title>{title}</Dialog.Title>
    //       {content}
    //     </Dialog.Panel>
    //     <Button onClick={onNo} variant="secondary">
    //       Cancel
    //     </Button>
    //     <Button
    //       onClick={onYes}
    //       variant={yesDestructive ? "destructive" : "indigo"}
    //     >
    //       {yesText}
    //     </Button>
    //   </Dialog>
    // </Transition>
    <Transition.Root show as={Fragment}>
      <Dialog as="div" className="relative z-[100]" onClose={noop}>
        <Transition.Child
          as={Fragment}
          enter="ease-out duration-300"
          enterFrom="opacity-0"
          enterTo="opacity-100"
          leave="ease-in duration-200"
          leaveFrom="opacity-100"
          leaveTo="opacity-0"
        >
          <div className="fixed inset-0 bg-gray-500 dark:bg-gray-900 dark:bg-opacity-90 bg-opacity-75 transition-opacity" />
        </Transition.Child>

        <div className="fixed inset-0 z-10 overflow-y-auto">
          <div className="flex min-h-full items-end justify-center p-4 text-center sm:items-center sm:p-0">
            <Transition.Child
              as={Fragment}
              enter="ease-out duration-300"
              enterFrom="opacity-0 translate-y-4 sm:translate-y-0 sm:scale-95"
              enterTo="opacity-100 translate-y-0 sm:scale-100"
              leave="ease-in duration-200"
              leaveFrom="opacity-100 translate-y-0 sm:scale-100"
              leaveTo="opacity-0 translate-y-4 sm:translate-y-0 sm:scale-95"
            >
              <Dialog.Panel className="relative transform overflow-hidden rounded-lg text-left shadow-xl transition-all sm:my-8 sm:w-full sm:max-w-lg text-gray-900 dark:text-gray-100">
                <div className="bg-white px-4 pt-5 pb-4 sm:p-6 sm:pb-4 dark:bg-gray-800">
                  <div className="sm:flex sm:items-start">
                    <div className="mt-3 text-center sm:mt-0 sm:ml-4 sm:text-left">
                      <Dialog.Title
                        as="h3"
                        className="text-lg font-medium leading-6"
                      >
                        {title}
                      </Dialog.Title>
                      <div className="mt-2">{content}</div>
                    </div>
                  </div>
                </div>
                <div className="bg-gray-50 px-4 py-3 sm:flex sm:flex-row-reverse sm:px-6 gap-4 dark:bg-gray-700">
                  <Button
                    onClick={onYes}
                    variant={yesDestructive ? "destructive" : "indigo"}
                  >
                    {yesText}
                  </Button>
                  <Button onClick={onNo} variant="transparent">
                    Cancel
                  </Button>
                </div>
              </Dialog.Panel>
            </Transition.Child>
          </div>
        </div>
      </Dialog>
    </Transition.Root>
  );
}

export function prompt(options: Omit<PromptStackItem, "id" | "onResolve">) {
  const id = new Array(16).join().replace(/(.|$)/g, function () {
    return ((Math.random() * 36) | 0).toString(36);
  });
  let resolveFunc: (v: boolean | null) => void = () => undefined;
  const promise = new Promise<boolean | null>((resolve) => {
    resolveFunc = resolve;
  });
  _promptStack.push({ ...options, onResolve: resolveFunc, id });
  return promise;
}
