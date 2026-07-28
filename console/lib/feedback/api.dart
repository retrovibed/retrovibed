import 'dart:convert';

import 'package:retrovibed/httpx.dart' as httpx;

// The GitHub App is installed once, on a single hardcoded repo. These node
// ids are stable for that repo/category and are not secrets - fetch them
// once via a GraphQL query against the repo and fill them in here.
//
//   query {
//     repository(owner: "retrovibed", name: "retrovibed") {
//       id
//       discussionCategories(first: 25) { nodes { id name } }
//     }
//   }
const String repositoryOwner = "retrovibed";
const String repositoryName = "retrovibed";
const String repositoryNodeId = "REPLACE_WITH_REPOSITORY_NODE_ID";
const String discussionCategoryNodeId = "REPLACE_WITH_DISCUSSION_CATEGORY_NODE_ID";

class GitHubInstallationToken {
  final String token;
  final DateTime expiresAt;

  const GitHubInstallationToken({required this.token, required this.expiresAt});

  factory GitHubInstallationToken.fromJson(Map<String, dynamic> json) {
    return GitHubInstallationToken(
      token: json["token"] as String,
      expiresAt: DateTime.parse(json["expires_at"] as String),
    );
  }

  bool get expired => DateTime.now().isAfter(expiresAt);
}

class GitHub {
  // fetches a short-lived GitHub App installation access token from deeppool.
  // deeppool holds the App's private key and derives eligibility from the
  // caller's customer_support permission (see meta.Token.customerSupport) -
  // this call never touches GitHub itself.
  static Future<GitHubInstallationToken> token({
    List<httpx.Option> options = const [],
  }) {
    return httpx
        .get(
          Uri.https(httpx.metaendpoint(), "/m/feedback/token"),
          options: [httpx.Accept.json, ...options],
        )
        .then((v) => GitHubInstallationToken.fromJson(jsonDecode(v.body) as Map<String, dynamic>));
  }

  static httpx.Option _githubBearer(String token) {
    return (httpx.Request request) {
      request.headers["Authorization"] = "Bearer $token";
      request.headers["Accept"] = "application/vnd.github+json";
      return Future.value(request);
    };
  }

  // creates an issue directly against GitHub's REST API using the installation
  // token obtained from GitHub.token(). Returns the created issue's html_url.
  static Future<String> createIssue({
    required String token,
    required String title,
    required String body,
  }) {
    return httpx
        .post(
          Uri.https("api.github.com", "/repos/$repositoryOwner/$repositoryName/issues"),
          options: [_githubBearer(token), httpx.Content.json],
          body: jsonEncode({"title": title, "body": body}),
        )
        .then((v) => (jsonDecode(v.body) as Map<String, dynamic>)["html_url"] as String);
  }

  // creates a discussion directly against GitHub's GraphQL API using the
  // installation token obtained from GitHub.token(). Returns the created
  // discussion's url.
  static Future<String> createDiscussion({
    required String token,
    required String title,
    required String body,
  }) {
    const mutation = r'''
      mutation($repositoryId: ID!, $categoryId: ID!, $title: String!, $body: String!) {
        createDiscussion(input: {repositoryId: $repositoryId, categoryId: $categoryId, title: $title, body: $body}) {
          discussion { url }
        }
      }
    ''';

    return httpx
        .post(
          Uri.https("api.github.com", "/graphql"),
          options: [_githubBearer(token), httpx.Content.json],
          body: jsonEncode({
            "query": mutation,
            "variables": {
              "repositoryId": repositoryNodeId,
              "categoryId": discussionCategoryNodeId,
              "title": title,
              "body": body,
            },
          }),
        )
        .then((v) {
          final decoded = jsonDecode(v.body) as Map<String, dynamic>;
          final data = decoded["data"] as Map<String, dynamic>;
          final created = data["createDiscussion"] as Map<String, dynamic>;
          final discussion = created["discussion"] as Map<String, dynamic>;
          return discussion["url"] as String;
        });
  }
}
