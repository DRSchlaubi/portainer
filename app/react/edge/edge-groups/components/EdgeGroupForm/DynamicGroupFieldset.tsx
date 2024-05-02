import { useFormikContext } from 'formik';

import { EdgeTypes } from '@/react/portainer/environments/types';
import { EdgeGroupAssociationTable } from '@/react/edge/components/EdgeGroupAssociationTable';

import { FormSection } from '@@/form-components/FormSection';
import { BoxSelector } from '@@/BoxSelector';
import { TagSelector } from '@@/TagSelector';

import { tagOptions } from './tag-options';
import { FormValues } from './types';

export function DynamicGroupFieldset() {
  const { values, setFieldValue, errors } = useFormikContext<FormValues>();
  return (
    <>
      <FormSection title="Tags">
        <BoxSelector
          slim
          value={values.partialMatch}
          onChange={(partialMatch) =>
            setFieldValue('partialMatch', partialMatch)
          }
          options={tagOptions}
          radioName="partialMatch"
        />
        <TagSelector
          value={values.tagIds}
          onChange={(tagIds) => setFieldValue('tagIds', tagIds)}
          errors={errors.tagIds}
        />
      </FormSection>

      <EdgeGroupAssociationTable
        data-cy="edgeGroupCreate-associatedEnvironmentsTable"
        title="Associated environments by tags"
        emptyContentLabel="No associated available"
        query={{
          types: EdgeTypes,
          tagIds: values.tagIds,
          tagsPartialMatch: values.partialMatch,
        }}
      />
    </>
  );
}
