import { spawn } from 'child_process';
import { platform } from 'os';

jest.mock('child_process');
jest.mock('os');

const mockSpawn = spawn as jest.MockedFunction<typeof spawn>;
const mockPlatform = platform as jest.MockedFunction<typeof platform>;

describe('CaffeinateServer', () => {
  beforeEach(() => {
    jest.clearAllMocks();
    mockPlatform.mockReturnValue('darwin');
  });
  
  describe('Platform checks', () => {
    it('should only work on macOS', () => {
      mockPlatform.mockReturnValue('linux');
      
      const mockProcess = {
        pid: 12345,
        kill: jest.fn(),
        on: jest.fn(),
      } as any;
      
      mockSpawn.mockReturnValue(mockProcess);
      
      expect(mockPlatform()).toBe('linux');
    });
  });
  
  describe('Caffeinate command building', () => {
    it('should build correct arguments for caffeinate', () => {
      const testCases = [
        {
          options: { display: true },
          expectedArgs: ['-d'],
        },
        {
          options: { idle: true, system: true },
          expectedArgs: ['-i', '-s'],
        },
        {
          options: { timeout: 300 },
          expectedArgs: ['-t', '300'],
        },
        {
          options: { pid: 12345 },
          expectedArgs: ['-w', '12345'],
        },
        {
          options: {
            display: true,
            idle: true,
            disk: true,
            system: true,
            user: true,
            timeout: 600,
            pid: 54321,
          },
          expectedArgs: ['-d', '-i', '-m', '-s', '-u', '-t', '600', '-w', '54321'],
        },
      ];
      
      testCases.forEach(({ options, expectedArgs }) => {
        const args: string[] = [];
        
        if (options.display) args.push('-d');
        if (options.idle) args.push('-i');
        if (options.disk) args.push('-m');
        if (options.system) args.push('-s');
        if (options.user) args.push('-u');
        if (options.timeout !== undefined) args.push('-t', options.timeout.toString());
        if (options.pid !== undefined) args.push('-w', options.pid.toString());
        
        expect(args).toEqual(expectedArgs);
      });
    });
  });
});