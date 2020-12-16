import React from 'react';
import { Prism as SyntaxHighlighter } from 'react-syntax-highlighter';

interface Props {
  value: string;
}

const Yaml: React.FC<Props> = (props: Props) => {
  const { value } = props;
  return (
    <SyntaxHighlighter language="yaml" showLineNumbers={true} wrapLines={true}>
      {value}
    </SyntaxHighlighter>
  );
};

export default Yaml;
