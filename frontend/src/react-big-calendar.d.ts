declare module "react-big-calendar" {
  import type { ComponentType, CSSProperties, ReactElement } from "react";

  export type View = "month" | "week" | "work_week" | "day" | "agenda";

  export interface Event {
    title?: string;
    start: Date;
    end: Date;
    resource?: unknown;
    allDay?: boolean;
  }

  /** Props passed to a custom `components.event` renderer */
  export interface EventProps<TEvent extends Event = Event> {
    event: TEvent;
  }

  export interface CalendarProps<TEvent extends Event = Event> {
    culture?: string;
    localizer: unknown;
    events: TEvent[];
    startAccessor: keyof TEvent | ((event: TEvent) => Date);
    endAccessor: keyof TEvent | ((event: TEvent) => Date);
    view?: View;
    date?: Date;
    onNavigate?: (date: Date) => void;
    onView?: (view: View) => void;
    selectable?: boolean;
    onSelectSlot?: (slotInfo: { start: Date; end: Date; slots: Date[]; action: string }) => void;
    views?: View[];
    style?: CSSProperties;
    eventPropGetter?: (event: TEvent) => { style?: CSSProperties; className?: string };
    components?: {
      toolbar?: ComponentType<Record<string, unknown>>;
      event?: ComponentType<EventProps<TEvent>>;
    };
    step?: number;
    timeslots?: number;
    min?: Date;
    max?: Date;
    scrollToTime?: Date;
    dayLayoutAlgorithm?: "overlap" | "no-overlap";
  }

  export const Calendar: <TEvent extends Event = Event>(props: CalendarProps<TEvent>) => ReactElement | null;

  // date-fns v4 typings are stricter than RBC’s localizer contract; keep loose here.
  export function dateFnsLocalizer(args: Record<string, unknown>): unknown;
}
