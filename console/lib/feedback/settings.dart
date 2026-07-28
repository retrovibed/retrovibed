import 'package:flutter/material.dart';
import 'package:retrovibed/authn.dart' as authn;
import 'package:retrovibed/designkit.dart' as ds;
import 'package:retrovibed/feedback/api.dart' as api;
import 'package:url_launcher/url_launcher.dart';

enum _SubmissionType { issue, discussion }

class Settings extends StatefulWidget {
  const Settings({super.key});

  @override
  State<Settings> createState() => _SettingsState();
}

class _SettingsState extends State<Settings> {
  bool _submitting = false;
  Widget _cause = ds.Error.zero;
  String? _createdUrl;
  _SubmissionType _type = _SubmissionType.issue;
  final _titleController = TextEditingController();
  final _bodyController = TextEditingController();

  void setState(VoidCallback fn) {
    if (!mounted) return;
    super.setState(fn);
  }

  @override
  void dispose() {
    _titleController.dispose();
    _bodyController.dispose();
    super.dispose();
  }

  Future<void> _submit() {
    setState(() {
      _submitting = true;
      _cause = ds.Error.zero;
      _createdUrl = null;
    });

    final title = _titleController.text;
    final body = _bodyController.text;

    return api.GitHub
        .token(options: [authn.request(authn.AuthzCache.meta(context))])
        .then((token) {
          switch (_type) {
            case _SubmissionType.issue:
              return api.GitHub.createIssue(token: token.token, title: title, body: body);
            case _SubmissionType.discussion:
              return api.GitHub.createDiscussion(token: token.token, title: title, body: body);
          }
        })
        .then((url) {
          setState(() {
            _submitting = false;
            _createdUrl = url;
          });
        })
        .catchError((e) {
          setState(() {
            _submitting = false;
            _cause = ds.Error.unknown(e, onTap: _submit);
          });
        });
  }

  @override
  Widget build(BuildContext context) {
    final theme = Theme.of(context);
    final defaults = ds.Defaults.of(context);

    return ds.Loading(
      loading: _submitting,
      cause: _cause,
      Padding(
        padding: defaults.padding,
        child: Column(
          crossAxisAlignment: CrossAxisAlignment.start,
          spacing: defaults.spacing,
          children: [
            Text("GitHub", style: theme.textTheme.titleMedium),
            SegmentedButton<_SubmissionType>(
              showSelectedIcon: false,
              expandedInsets: EdgeInsets.zero,
              segments: const [
                ButtonSegment(value: _SubmissionType.issue, label: Text("Issue")),
                ButtonSegment(value: _SubmissionType.discussion, label: Text("Discussion")),
              ],
              selected: {_type},
              onSelectionChanged: (selection) => setState(() => _type = selection.first),
            ),
            TextFormField(
              controller: _titleController,
              decoration: const InputDecoration(helperText: "title"),
              autofocus: true,
            ),
            TextFormField(
              controller: _bodyController,
              decoration: const InputDecoration(helperText: "description"),
              minLines: 4,
              maxLines: 12,
            ),
            if (_createdUrl != null)
              InkWell(
                onTap: () => launchUrl(Uri.parse(_createdUrl!)),
                child: Text(
                  _createdUrl!,
                  style: theme.textTheme.bodySmall?.copyWith(
                    color: theme.colorScheme.primary,
                    decoration: TextDecoration.underline,
                  ),
                ),
              ),
            SizedBox(
              width: double.infinity,
              child: ds.LoadingButton(
                Text(_type == _SubmissionType.issue ? "Create Issue" : "Start Discussion"),
                onPressed: _submit,
              ),
            ),
          ],
        ),
      ),
    );
  }
}
