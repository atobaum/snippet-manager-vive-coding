# README: Snippet Management Tool "sni"

## 1. 목적 (Purpose) 🎯

`sni`는 개발자의 파편화된 코드 조각, 명령어, 설정 파일 등을 **중앙에서 관리**하고, **가장 익숙한 환경인 터미널(CLI)과 브라우저(Web UI)**를 통해 빠르고 쉽게 접근할 수 있도록 돕는 도구입니다. 이를 통해 컨텍스트 전환 비용을 줄이고 개발 생산성을 극대화하는 것을 목표로 합니다.

---

## 2. 핵심 기능 (Core Features) ✨

`sni`는 커맨드 라인과 웹 UI, 두 가지 인터페이스를 통해 일관된 스니펫 관리 경험을 제공합니다.

### CLI (Command-Line Interface)

* **`sni new <name>`**: 새로운 스니펫을 등록합니다.
* **`sni edit <name>`**: 기존 스니펫을 수정합니다.
* **`sni list [--color]`**: 저장된 모든 스니펫의 목록을 간략히 보여줍니다.
* **`sni search <keyword> [--color]`**: 키워드로 스니펫을 검색합니다.
* **`sni use <name>`**: 스니펫의 내용을 터미널에 출력하여 바로 사용하거나 다른 명령어와 조합할 수 있습니다.
* **`sni exec [--tag <tag>] [--color]`**: 🆕 인터랙티브하게 스니펫을 선택하고 실행합니다 (fzf 지원).
* **`sni rm <name>`**: 스니펫을 삭제합니다.
* **`sni configure`**: 🆕 설정 정보를 확인합니다.
* **`sni server [--dev] [--port <port>]`**: 스니펫 관리를 위한 로컬 웹 UI를 실행합니다.

### Web UI (Graphical User Interface)

* **대시보드**: 모든 스니펫을 한눈에 볼 수 있는 시각적인 대시보드를 제공합니다.
* **실시간 검색**: 키워드를 입력하여 스니펫을 빠르게 필터링합니다.
* **간편한 관리**: 웹 폼(Form)을 통해 스니펫을 직관적으로 생성, 수정, 삭제할 수 있습니다.
* **클립보드 복사**: 클릭 한 번으로 스니펫의 내용을 클립보드에 복사할 수 있습니다.

---

## 3. 기술 스택 (Tech Stack) 🛠️

- backend: go, cobra, embed
- frontend: svelte static site

go backend server가 svelte static site를 제공하는 방식으로 구현됩니다.

---

## 4. 데이터 모델 (Data Model) 📄

모든 스니펫은 **`.sni/snippets.yaml` 단일 파일**에 저장되어 관리가 용이합니다.

* **파일 구조 예시:**

    ```yaml
    # .sni/snippets.yaml
    snippets:
      k8s-pod:
        description: "Nginx Pod를 생성하는 기본 매니페스트"
        tags: ["k8s", "pod", "nginx"]
        command: |
          apiVersion: v1
          kind: Pod
          metadata:
            name: nginx-pod
          spec:
            containers:
            - name: nginx
              image: nginx:1.14.2
              ports:
              - containerPort: 80
      
      find-large-files:
        description: "특정 디렉토리에서 용량이 큰 파일을 찾는 셸 스크립트"
        tags: ["shell", "utility"]
        command: |
          find . -type f -size +100M -exec ls -lh {} \; | awk '{ print $9 ": " $5 }'
    ```

---

## 5. 새로운 기능 (New Features) ✨

### 🎯 인터랙티브 실행 (exec 명령어)
- **fzf 통합**: fzf가 설치된 경우 fuzzy finder로 스니펫 선택
- **fallback 지원**: fzf가 없어도 번호 기반 선택으로 동작
- **태그 필터링**: `--tag` 옵션으로 특정 태그의 스니펫만 표시
- **실행 확인**: 실행 전 명령어 내용 확인 및 승인

### 🎨 컬러 출력
- **`--color` 플래그**: list, search, exec 명령어에서 컬러화된 출력
- **구문 강조**: 스니펫 이름, 설명, 태그, 명령어를 다른 색상으로 표시
- **상태 메시지**: 성공, 오류, 경고 메시지를 색상으로 구분

### ⚙️ 개발 모드
- **`--dev` 플래그**: Svelte 개발 서버와 연동하여 실시간 개발
- **프록시 지원**: API 요청은 Go 서버로, 프론트엔드는 Vite 서버로 자동 프록시
- **핫 리로드**: 프론트엔드 코드 변경 시 자동 새로고침

---

## 6. 설치 및 실행 (Installation & Usage) 🛠️

### 빌드

```bash
# 의존성 설치 및 빌드
./build.sh
```

### CLI 사용법

```bash
# 도움말 보기
./sni --help

# 새 스니펫 생성
./sni new my-snippet

# 모든 스니펫 목록 보기
./sni list
./sni list --color          # 컬러 출력

# 키워드로 스니펫 검색
./sni search docker
./sni search docker --color # 컬러 출력

# 스니펫 내용 출력
./sni use my-snippet

# 스니펫 수정
./sni edit my-snippet

# 스니펫 삭제
./sni rm my-snippet

# 스니펫 실행 (인터랙티브)
./sni exec                  # fzf 또는 번호 선택으로 실행
./sni exec --tag docker     # 태그로 필터링
./sni exec --color          # 컬러 출력

# 설정 확인
./sni configure

# 웹 UI 서버 시작
./sni server
./sni server --dev          # 개발 모드 (Svelte dev server와 연동)
./sni server --port 9090    # 커스텀 포트
```

### 웹 UI

```bash
# 웹 서버 시작
./sni server

# 브라우저에서 http://localhost:8080 접속
```

웹 UI에서는 다음 기능을 제공합니다:
- 📝 스니펫 생성, 수정, 삭제
- 🔍 실시간 검색 및 필터링
- 📋 클립보드로 복사
- 🏷️ 태그 기반 분류
- 📱 반응형 디자인

---

## 7. 설정 (Configuration) ⚙️

환경변수를 통해 설정 디렉토리를 변경할 수 있습니다:

```bash
# 커스텀 설정 디렉토리 사용
export SNI_CONFIG_DIR="/path/to/your/config"
./sni list
```

기본값:
- 현재 디렉토리: `.sni/snippets.yaml`
- 홈 디렉토리: `~/.config/sni/snippets.yaml` (fallback)