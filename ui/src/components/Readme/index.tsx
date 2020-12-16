import React from 'react';
import { Prism as SyntaxHighlighter } from 'react-syntax-highlighter';

interface Props {
  value: string;
}

const Readme: React.FC<Props> = (props: Props) => {
  const { value } = props;
  return (
    <SyntaxHighlighter language="markdown" showLineNumbers={true} wrapLines={true}>
      {value}
    </SyntaxHighlighter>
  );
};

export default Readme;
