# README: Snippet Management Tool "sni"

## 1. 목적 (Purpose) 🎯

`sni`는 개발자의 파편화된 코드 조각, 명령어, 설정 파일 등을 **중앙에서 관리**하고, **가장 익숙한 환경인 터미널(CLI)과 브라우저(Web UI)**를 통해 빠르고 쉽게 접근할 수 있도록 돕는 도구입니다. 이를 통해 컨텍스트 전환 비용을 줄이고 개발 생산성을 극대화하는 것을 목표로 합니다.

---

## 2. 핵심 기능 (Core Features) ✨

`sni`는 커맨드 라인과 웹 UI, 두 가지 인터페이스를 통해 일관된 스니펫 관리 경험을 제공합니다.

### CLI (Command-Line Interface)

* **`sni new <name>`**: 새로운 스니펫을 등록합니다.
* **`sni edit <name>`**: 기존 스니펫을 수정합니다.
* **`sni list`**: 저장된 모든 스니펫의 목록을 간략히 보여줍니다.
* **`sni search <keyword>`**: 키워드로 스니펫을 검색합니다.
* **`sni use <name>`**: 스니펫의 내용을 터미널에 출력하여 바로 사용하거나 다른 명령어와 조합할 수 있습니다.
* **`sni rm <name>`**: 스니펫을 삭제합니다.
* **`sni server`**: 스니펫 관리를 위한 로컬 웹 UI를 실행합니다.

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

## 5. 설치 및 실행 (Installation & Usage) 🛠️

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

# 키워드로 스니펫 검색
./sni search docker

# 스니펫 내용 출력
./sni use my-snippet

# 스니펫 수정
./sni edit my-snippet

# 스니펫 삭제
./sni rm my-snippet

# 웹 UI 서버 시작
./sni server
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

## 6. 설정 (Configuration) ⚙️

환경변수를 통해 설정 디렉토리를 변경할 수 있습니다:

```bash
# 커스텀 설정 디렉토리 사용
export SNI_CONFIG_DIR="/path/to/your/config"
./sni list
```

기본값:
- 현재 디렉토리: `.sni/snippets.yaml`
- 홈 디렉토리: `~/.config/sni/snippets.yaml` (fallback)