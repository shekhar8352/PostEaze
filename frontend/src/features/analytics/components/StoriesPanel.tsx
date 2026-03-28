import { Paper, Text, Table, Skeleton } from "@mantine/core";
import type { StoryAnalyticsItem } from "../types";
import styles from "./StoriesPanel.module.css";

type Props = {
  stories: StoryAnalyticsItem[];
  loading?: boolean;
};

function cell(n: number | undefined) {
  if (n == null) return "—";
  return n.toLocaleString();
}

export function StoriesPanel({ stories, loading }: Props) {
  if (!loading && stories.length === 0) {
    return null;
  }

  return (
    <Paper className={styles.wrap} p="md" radius="md" withBorder>
      <Text fw={600} size="sm" mb="md" c="var(--pe-text)">
        Story performance
      </Text>
      {loading ? (
        <Skeleton height={200} />
      ) : (
        <Table.ScrollContainer minWidth={640}>
          <Table verticalSpacing="sm" striped highlightOnHover>
            <Table.Thead>
              <Table.Tr>
                <Table.Th>Date</Table.Th>
                <Table.Th>Caption</Table.Th>
                <Table.Th>Impressions</Table.Th>
                <Table.Th>Reach</Table.Th>
                <Table.Th>Taps fwd</Table.Th>
                <Table.Th>Taps back</Table.Th>
                <Table.Th>Exits</Table.Th>
                <Table.Th>Replies</Table.Th>
              </Table.Tr>
            </Table.Thead>
            <Table.Tbody>
              {stories.map((s) => (
                <Table.Tr key={`${s.id}-${s.date}`}>
                  <Table.Td>{s.date}</Table.Td>
                  <Table.Td>
                    <Text size="sm" lineClamp={1}>
                      {s.caption || "—"}
                    </Text>
                  </Table.Td>
                  <Table.Td>{cell(s.impressions)}</Table.Td>
                  <Table.Td>{cell(s.reach)}</Table.Td>
                  <Table.Td>{cell(s.taps_forward)}</Table.Td>
                  <Table.Td>{cell(s.taps_backward)}</Table.Td>
                  <Table.Td>{cell(s.exits)}</Table.Td>
                  <Table.Td>{cell(s.replies)}</Table.Td>
                </Table.Tr>
              ))}
            </Table.Tbody>
          </Table>
        </Table.ScrollContainer>
      )}
    </Paper>
  );
}
