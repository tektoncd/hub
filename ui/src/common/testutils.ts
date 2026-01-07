const RealDate = global.Date;
const FIXED = new RealDate('2026-01-01T10:20:30Z');

export const FakeDate = () => {
  beforeEach(() => {
    global.Date = class extends RealDate {
      constructor(...args: [number?]) {
        super(args[0] ?? FIXED.getTime());
      }
      static now = () => FIXED.getTime();
    } as DateConstructor;
  });
  afterEach(() => {
    global.Date = RealDate;
  });
};

export const ActualDate = () => {};
