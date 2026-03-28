import { useMemo, useState } from "react";
import {
  Paper,
  Text,
  Table,
  Skeleton,
  Group,
  Select,
  Badge,
  Anchor,
  Box,
  UnstyledButton,
} from "@mantine/core";
import { Icons } from "@/app/theme";
import type { TopPostItem, TopPostsSort } from "../types";
import styles from "./TopPostsTable.module.css";

type SortKey = "engagement" | "reach" | "impressions" | "likes" | "plays";

type Props = {
  posts: TopPostItem[];
  loading?: boolean;
  apiSort: TopPostsSort;
  onApiSortChange: (s: TopPostsSort) => void;
  postType: string;
  onPostTypeChange: (v: string) => void;
};

function truncate(s: string, n: number) {
  const t = s.replace(/\s+/g, " ").trim();
  if (t.length <= n) return t || "—";
  return `${t.slice(0, n)}…`;
}

function Th({ label, active, reversed, onClick }: { label: string; active: boolean; reversed?: boolean; onClick: () => void }) {
  return (
    <Table.Th className={styles.th}>
      <UnstyledButton onClick={onClick} className={styles.thBtn}>
        <Group gap={4} wrap="nowrap">
          <Text size="xs" fw={600}>
            {label}
          </Text>
          {active ? (
            reversed ? (
              <Icons.ChevronUp size={14} />
            ) : (
              <Icons.ChevronDown size={14} />
            )
          ) : null}
        </Group>
      </UnstyledButton>
    </Table.Th>
  );
}

export function TopPostsTable({
  posts,
  loading,
  apiSort,
  onApiSortChange,
  postType,
  onPostTypeChange,
}: Props) {
  const [sortKey, setSortKey] = useState<SortKey>("engagement");
  const [reverse, setReverse] = useState(true);

  const sorted = useMemo(() => {
    const arr = [...posts];
    const mult = reverse ? -1 : 1;
    arr.sort((a, b) => {
      const va = a[sortKey] ?? 0;
      const vb = b[sortKey] ?? 0;
      return (va - vb) * mult;
    });
    return arr;
  }, [posts, sortKey, reverse]);

  const toggle = (key: SortKey) => {
    if (sortKey === key) setReverse(!reverse);
    else {
      setSortKey(key);
      setReverse(true);
    }
  };

  return (
    <Paper className={styles.wrap} p="md" radius="md" withBorder>
      <Group justify="space-between" align="flex-end" mb="md" wrap="wrap" gap="sm">
        <Text fw={600} size="sm" c="var(--pe-text)">
          Top posts
        </Text>
        <Group gap="sm">
          <Select
            label="Rank by"
            size="xs"
            w={140}
            value={apiSort}
            onChange={(v) => v && onApiSortChange(v as TopPostsSort)}
            data={[
              { value: "engagement", label: "Engagement" },
              { value: "reach", label: "Reach" },
              { value: "impressions", label: "Impressions" },
              { value: "plays", label: "Plays" },
            ]}
          />
          <Select
            label="Post type"
            size="xs"
            w={120}
            value={postType}
            onChange={(v) => onPostTypeChange(v ?? "")}
            data={[
              { value: "", label: "All" },
              { value: "image", label: "Image" },
              { value: "video", label: "Video" },
              { value: "reel", label: "Reel" },
              { value: "carousel", label: "Carousel" },
              { value: "story", label: "Story" },
            ]}
          />
        </Group>
      </Group>

      {loading ? (
        <Skeleton height={240} />
      ) : sorted.length === 0 ? (
        <Text c="dimmed" size="sm">
          No posts with analytics in this range.
        </Text>
      ) : (
        <Table.ScrollContainer minWidth={720}>
          <Table verticalSpacing="sm" striped highlightOnHover>
            <Table.Thead>
              <Table.Tr>
                <Table.Th>Preview</Table.Th>
                <Table.Th>Caption</Table.Th>
                <Table.Th>Type</Table.Th>
                <Th label="Reach" active={sortKey === "reach"} reversed={reverse} onClick={() => toggle("reach")} />
                <Th label="Impr." active={sortKey === "impressions"} reversed={reverse} onClick={() => toggle("impressions")} />
                <Th label="Likes" active={sortKey === "likes"} reversed={reverse} onClick={() => toggle("likes")} />
                <Th label="Engagement" active={sortKey === "engagement"} reversed={reverse} onClick={() => toggle("engagement")} />
                <Th label="Plays" active={sortKey === "plays"} reversed={reverse} onClick={() => toggle("plays")} />
              </Table.Tr>
            </Table.Thead>
            <Table.Tbody>
              {sorted.map((row) => (
                <Table.Tr key={row.post_id}>
                  <Table.Td>
                    <Box className={styles.thumbWrap}>
                      {row.thumbnail_url ? (
                        <img src={row.thumbnail_url} alt="" className={styles.thumb} loading="lazy" />
                      ) : (
                        <div className={styles.thumbPlaceholder}>
                          <Icons.Photo size={20} />
                        </div>
                      )}
                    </Box>
                  </Table.Td>
                  <Table.Td>
                    <Text size="sm" lineClamp={2}>
                      {truncate(row.caption, 120)}
                    </Text>
                    {row.permalink ? (
                      <Anchor href={row.permalink} size="xs" target="_blank" rel="noreferrer">
                        Open on Instagram
                      </Anchor>
                    ) : null}
                  </Table.Td>
                  <Table.Td>
                    <Badge size="sm" variant="light">
                      {row.post_type}
                    </Badge>
                  </Table.Td>
                  <Table.Td>{row.reach.toLocaleString()}</Table.Td>
                  <Table.Td>{row.impressions.toLocaleString()}</Table.Td>
                  <Table.Td>{row.likes.toLocaleString()}</Table.Td>
                  <Table.Td>{row.engagement.toLocaleString()}</Table.Td>
                  <Table.Td>{row.plays.toLocaleString()}</Table.Td>
                </Table.Tr>
              ))}
            </Table.Tbody>
          </Table>
        </Table.ScrollContainer>
      )}
    </Paper>
  );
}
