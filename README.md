# Drift 🫘
*Dynamic Runtime Integration & Fast Testing*

Drift is an agnostic, local-first mobile development CLI tool that automates building and deployment for any mobile project. Think EAS but framework-agnostic and starting with local builds. Built with Go for speed and reliability.

## 🎯 Vision

Make mobile app development and deployment as fast and simple as web development. Start with a powerful local CLI tool, then evolve into a full platform. Don't reinvent the wheel - leverage existing tools like Fastlane, Gradle, and Xcode CLI.

## ✨ Core Features

### 🚀 Build & Deploy (Phase 1 - Local First)
- **Universal Building**: Support for React Native, Flutter, Ionic, Cordova, and native iOS/Android projects
- **Local Builds**: Fast local builds using existing toolchains (Xcode, Gradle, Flutter CLI)
- **Automated Store Submission**: Direct uploads using Fastlane and existing store APIs
- **TestFlight/Internal Testing**: Automatic beta distribution via Fastlane
- **Certificate Management**: Leverage Fastlane's cert and sigh tools
- **Multi-Environment Support**: Dev, staging, production configs

### 🔧 Development Tools (Phase 2)
- **Native Log Viewer**: Real-time device logs using `adb logcat` and `xcrun simctl`
- **Build Status**: Clear build progress and error reporting
- **Device Management**: List and target connected devices/simulators
- **Configuration Validation**: Validate project settings and dependencies

### 🛠 Future Platform Features (Phase 3+)
- **Web Dashboard**: Manage multiple projects and builds
- **Cloud Builds**: Remote building without local dependencies  
- **Team Collaboration**: Share builds and manage access
- **Advanced Analytics**: Build times, success rates, deployment tracking
- **Network Inspector**: HTTP request monitoring and debugging
- **Performance Monitoring**: Memory usage, crash reports, performance metrics

### 🤖 Automation & CI/CD (Phase 4)
- **Git Integration**: Automatic builds on commits, tags, or PR merges
- **Playwright Store Automation**: Advanced App Store Connect and Google Play Console automation
- **Notifications**: Slack/Discord/Email build status updates
- **Rollback Support**: Quick rollback to previous versions
- **A/B Testing**: Percentage-based rollouts

## 🗺 Development Roadmap

### Phase 1: Local CLI Foundation (MVP)
- [ ] Go CLI setup with Cobra framework
- [ ] Project detection (React Native, Flutter, Ionic, Cordova, Native)
- [ ] Configuration file support (`drift.yml` or `drift.json`)
- [ ] iOS local builds (wrapping `xcodebuild` and Fastlane)
- [ ] Android local builds (wrapping Gradle and Fastlane)
- [ ] Device/simulator deployment
- [ ] Basic error handling and logging

### Phase 2: Store Automation & Dev Tools  
- [ ] Fastlane integration for App Store uploads
- [ ] Fastlane integration for Google Play uploads
- [ ] TestFlight automation
- [ ] Certificate and provisioning profile management (via Fastlane)
- [ ] Real-time device logs (`adb logcat`, iOS Console)
- [ ] Build artifact management
- [ ] Multi-environment configurations

### Phase 3: Platform Evolution
- [ ] Web dashboard (Go backend + React/Vue frontend)
- [ ] Build history and analytics
- [ ] Cloud build infrastructure
- [ ] Team collaboration features
- [ ] Advanced debugging tools

### Phase 4: Advanced Automation
- [ ] CI/CD pipeline integrations
- [ ] Playwright-based store automation
- [ ] Advanced notifications and webhooks
- [ ] A/B testing and gradual rollouts

## 🚀 Getting Started

### Prerequisites
- Go 1.21+
- Git
- Xcode (for iOS builds)
- Android Studio or Android SDK (for Android builds)

### Installation
```bash
go install github.com/your-org/drift@latest
# or download binary from releases
curl -sSL https://install.drift.dev | sh
```

### Quick Start
```bash
# Initialize drift in your mobile project
drift init

# Build and deploy to staging
drift deploy --env staging

# Deploy to production
drift deploy --env production --platform ios,android

# Start development tools
drift dev --logs --network-inspector
```

### Configuration Example
```yaml
# drift.yml
name: MyAwesomeApp
platforms: [ios, android]

ios:
  project_path: ./ios/MyApp.xcworkspace
  scheme: MyApp
  bundle_id: com.company.myapp
  
android:
  project_path: ./android
  module: app
  application_id: com.company.myapp

environments:
  staging:
    ios:
      bundle_id: com.company.myapp.staging
      provisioning_profile: "MyApp Staging"
      export_method: app-store
    android:
      application_id: com.company.myapp.staging
      build_type: debug
      
  production:
    ios:
      bundle_id: com.company.myapp
      provisioning_profile: "MyApp Production" 
      export_method: app-store
    android:
      application_id: com.company.myapp
      build_type: release

fastlane:
  ios_lane: beta
  android_lane: beta
  
testers:
  internal:
    - developer@company.com
    - tester@company.com
```

## 🛠 Development Setup

### Step-by-Step Implementation Guide

#### 1. Project Setup
```bash
mkdir drift
cd drift
go mod init github.com/your-org/drift

# Install dependencies
go get github.com/spf13/cobra@latest
go get github.com/spf13/viper@latest
go get gopkg.in/yaml.v3
go get github.com/fatih/color
go get github.com/briandowns/spinner
```

#### 2. Project Structure
```
├── cmd/
│   ├── root.go
│   ├── init.go
│   ├── build.go
│   ├── deploy.go
│   ├── logs.go
│   └── devices.go
├── internal/
│   ├── config/
│   │   ├── config.go
│   │   └── validate.go
│   ├── detect/
│   │   ├── project.go
│   │   └── platforms.go
│   ├── build/
│   │   ├── ios.go
│   │   ├── android.go
│   │   ├── flutter.go
│   │   └── react_native.go
│   ├── deploy/
│   │   ├── fastlane.go
│   │   ├── appstore.go
│   │   └── playstore.go
│   └── utils/
│       ├── exec.go
│       ├── logger.go
│       └── files.go
├── pkg/
│   └── drift/
└── main.go
```

#### 3. Technology Stack
- **CLI Framework**: Cobra (de facto standard for Go CLIs)
- **Configuration**: Viper + YAML/JSON support
- **Build Automation**: Wrap existing tools (xcodebuild, gradle, flutter CLI)
- **Store Deployment**: Fastlane integration
- **Process Management**: os/exec with proper error handling
- **Logging**: Structured logging with color output

#### 4. Core Components to Build First
1. **Project Detector**: Identify project type by scanning files
2. **Config Parser**: Handle `drift.yml` loading and validation  
3. **Build Orchestrator**: Coordinate platform-specific builds
4. **Fastlane Wrapper**: Execute Fastlane lanes with proper error handling

#### 5. First Commands to Implement
```bash
# Initialize drift config
drift init

# Detect project type and show info
drift info

# Build for specific platforms
drift build --platform ios
drift build --platform android
drift build --all

# Deploy using Fastlane
drift deploy --env staging
drift deploy --env production --platform ios

# Show connected devices
drift devices

# View real-time logs
drift logs --device "iPhone 14" --platform ios
```

#### 6. MVP Implementation Order
1. **CLI Structure**: Basic Cobra setup with commands
2. **Project Detection**: Identify RN, Flutter, Ionic, Native projects
3. **Configuration**: Parse drift.yml and validate settings
4. **iOS Build**: Wrap xcodebuild commands
5. **Android Build**: Wrap gradle commands
6. **Fastlane Integration**: Execute lanes for deployment
7. **Device Management**: List and target devices/simulators

## 🤝 Contributing

We welcome contributions! This is an ambitious project that could revolutionize mobile development workflow.

### Areas needing help:
- Fastlane automation scripts
- Playwright store automation
- Native log parsing
- Network request interception
- Cloud infrastructure setup
- Documentation and examples

## 📄 License

MIT License - see LICENSE file for details

---

*Built with ❤️ to make mobile development as enjoyable as web development*
