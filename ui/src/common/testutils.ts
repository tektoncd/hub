//Sets a Fake date for test
export const FakeDate = () => {
  beforeEach(() => {
    global.Date.now = jest.fn(() => new Date('2021-01-01T10:20:30Z').getTime());
  });
};

//Sets an Actual date for test
export const ActualDate = () => {
  afterEach(() => {
    global.Date = Date;
  });
};
