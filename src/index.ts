#!/usr/bin/env node
import { Server } from '@modelcontextprotocol/sdk/server/index.js';
import { StdioServerTransport } from '@modelcontextprotocol/sdk/server/stdio.js';
import {
  CallToolRequestSchema,
  ListToolsRequestSchema,
  Tool,
} from '@modelcontextprotocol/sdk/types.js';
import { spawn, ChildProcess } from 'child_process';
import { platform } from 'os';

interface CaffeinateProcess {
  id: string;
  process: ChildProcess;
  startTime: Date;
  options: CaffeinateOptions;
}

interface CaffeinateOptions {
  display?: boolean;
  idle?: boolean;
  disk?: boolean;
  system?: boolean;
  user?: boolean;
  timeout?: number;
  pid?: number;
}

class CaffeinateServer {
  private server: Server;
  private processes: Map<string, CaffeinateProcess> = new Map();
  
  constructor() {
    this.server = new Server(
      {
        name: 'caffeinate-mcp',
        version: '1.0.0',
      },
      {
        capabilities: {
          tools: {},
        },
      }
    );
    
    this.setupHandlers();
  }
  
  private setupHandlers() {
    this.server.setRequestHandler(ListToolsRequestSchema, async () => ({
      tools: this.getTools(),
    }));
    
    this.server.setRequestHandler(CallToolRequestSchema, async (request) => {
      const { name, arguments: args } = request.params;
      
      switch (name) {
        case 'caffeinate_start':
          return this.startCaffeinate(args as CaffeinateOptions);
        case 'caffeinate_stop':
          return this.stopCaffeinate((args as { id: string }).id);
        case 'caffeinate_list':
          return this.listCaffeinate();
        default:
          throw new Error(`Unknown tool: ${name}`);
      }
    });
  }
  
  private getTools(): Tool[] {
    return [
      {
        name: 'caffeinate_start',
        description: 'Start caffeinate to prevent system sleep',
        inputSchema: {
          type: 'object',
          properties: {
            display: {
              type: 'boolean',
              description: 'Prevent display from sleeping (-d flag)',
            },
            idle: {
              type: 'boolean',
              description: 'Prevent system from idle sleeping (-i flag)',
            },
            disk: {
              type: 'boolean',
              description: 'Prevent disk from idle sleeping (-m flag)',
            },
            system: {
              type: 'boolean',
              description: 'Prevent system from sleeping when on AC power (-s flag)',
            },
            user: {
              type: 'boolean',
              description: 'Declare user is active (-u flag)',
            },
            timeout: {
              type: 'number',
              description: 'Timeout in seconds (-t flag)',
            },
            pid: {
              type: 'number',
              description: 'Wait for process with specified PID to exit (-w flag)',
            },
          },
        },
      },
      {
        name: 'caffeinate_stop',
        description: 'Stop a caffeinate process',
        inputSchema: {
          type: 'object',
          properties: {
            id: {
              type: 'string',
              description: 'ID of the caffeinate process to stop',
            },
          },
          required: ['id'],
        },
      },
      {
        name: 'caffeinate_list',
        description: 'List active caffeinate processes',
        inputSchema: {
          type: 'object',
          properties: {},
        },
      },
    ];
  }
  
  private async startCaffeinate(options: CaffeinateOptions) {
    if (platform() !== 'darwin') {
      return {
        content: [
          {
            type: 'text',
            text: 'Error: caffeinate is only available on macOS',
          },
        ],
      };
    }
    
    const args: string[] = [];
    
    if (options.display) args.push('-d');
    if (options.idle) args.push('-i');
    if (options.disk) args.push('-m');
    if (options.system) args.push('-s');
    if (options.user) args.push('-u');
    if (options.timeout !== undefined) args.push('-t', options.timeout.toString());
    if (options.pid !== undefined) args.push('-w', options.pid.toString());
    
    try {
      const process = spawn('caffeinate', args);
      const id = `${process.pid}_${Date.now()}`;
      
      const caffeinateProcess: CaffeinateProcess = {
        id,
        process,
        startTime: new Date(),
        options,
      };
      
      this.processes.set(id, caffeinateProcess);
      
      process.on('exit', () => {
        this.processes.delete(id);
      });
      
      process.on('error', () => {
        this.processes.delete(id);
      });
      
      return {
        content: [
          {
            type: 'text',
            text: `Started caffeinate process with ID: ${id}`,
          },
        ],
      };
    } catch (error) {
      return {
        content: [
          {
            type: 'text',
            text: `Error starting caffeinate: ${error instanceof Error ? error.message : 'Unknown error'}`,
          },
        ],
      };
    }
  }
  
  private async stopCaffeinate(id: string) {
    const caffeinateProcess = this.processes.get(id);
    
    if (!caffeinateProcess) {
      return {
        content: [
          {
            type: 'text',
            text: `Error: No caffeinate process found with ID: ${id}`,
          },
        ],
      };
    }
    
    try {
      caffeinateProcess.process.kill();
      this.processes.delete(id);
      
      return {
        content: [
          {
            type: 'text',
            text: `Stopped caffeinate process with ID: ${id}`,
          },
        ],
      };
    } catch (error) {
      return {
        content: [
          {
            type: 'text',
            text: `Error stopping caffeinate: ${error instanceof Error ? error.message : 'Unknown error'}`,
          },
        ],
      };
    }
  }
  
  private async listCaffeinate() {
    const processList = Array.from(this.processes.entries()).map(([id, proc]) => {
      const flags = [];
      if (proc.options.display) flags.push('-d');
      if (proc.options.idle) flags.push('-i');
      if (proc.options.disk) flags.push('-m');
      if (proc.options.system) flags.push('-s');
      if (proc.options.user) flags.push('-u');
      if (proc.options.timeout !== undefined) flags.push(`-t ${proc.options.timeout}`);
      if (proc.options.pid !== undefined) flags.push(`-w ${proc.options.pid}`);
      
      return {
        id,
        pid: proc.process.pid,
        startTime: proc.startTime.toISOString(),
        flags: flags.join(' '),
      };
    });
    
    return {
      content: [
        {
          type: 'text',
          text: processList.length > 0
            ? JSON.stringify(processList, null, 2)
            : 'No active caffeinate processes',
        },
      ],
    };
  }
  
  async run() {
    const transport = new StdioServerTransport();
    await this.server.connect(transport);
  }
}

async function main() {
  const server = new CaffeinateServer();
  await server.run();
}

main().catch((error) => {
  console.error('Server error:', error);
  process.exit(1);
});