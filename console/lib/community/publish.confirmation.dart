import 'package:flutter/material.dart';
import 'package:retrovibed/designkit.dart' as ds;
import 'package:retrovibed/design.kit/forms.dart' as forms;
import 'package:retrovibed/uuidx.dart' as uuidx;
import 'package:retrovibed/httpx.dart' as httpx;
import 'package:retrovibed/authn.dart' as authn;
import 'package:retrovibed/media/media.pb.dart';
import 'package:retrovibed/media/media.known.pb.dart';
import 'package:retrovibed/mimex.dart' as mimex;
import 'api.dart';
import 'publish.mode.edit.dart';
import 'package:retrovibed/google/api.dart' as google;
import 'package:retrovibed/design.kit/stateful.dart';

class PublishConfirmation extends StatefulWidget {
  final Download? download;
  final Community? community;
  final Known? knownMedia;
  final VoidCallback onPublished;
  final Future<YouTubeStatus> Function({List<httpx.Option> options}) youtubeStatus;
  final Future<PublishContentResponse> Function(
    String cid,
    PublishContentRequest req, {
    List<httpx.Option> options,
  })
  apicommunitypublish;

  const PublishConfirmation({
    super.key,
    required this.download,
    required this.community,
    this.knownMedia,
    required this.onPublished,
    this.youtubeStatus = google.YouTube.status,
    this.apicommunitypublish = publishing.publish,
  });

  @override
  State<PublishConfirmation> createState() => _PublishConfirmationState();
}

class _PublishConfirmationState extends State<PublishConfirmation> with LoadingState {
  late PublishContentRequest _request;
  String _oauthGoogleId = '';
  bool _loading = false;
  Widget _cause = ds.Error.zero;

  @override
  void initState() {
    super.initState();
    _request = PublishContentRequest(
      publishMode: widget.community?.defaultPublishMode ?? PublishMode.UNLISTED,
      publishedContent: PublishedContent(),
    );
    ds.postframe(() {
      httpx
          .withRetry(
            () => widget.youtubeStatus(
              options: [authn.request(authn.AuthzCache.meta(context))],
            ),
          )
          .then((status) {
            setState(() {
              _oauthGoogleId = status.id;
            });
          })
          .catchError((e) {
            setState(() {
              _cause = ds.Errors.httpauto(e, onTap: _reseterr);
            });
          }, test: httpx.ErrorsTest.httpauto);
    });
  }

  void _reseterr() {
    setState(() {
      _cause = ds.Error.zero;
    });
  }

  void _publish() {
    setState(() {
      _loading = true;
      _cause = ds.Error.zero;
    });

    _request.publishedContent
      ..communityId = widget.community?.id ?? uuidx.min()
      ..knownMediaId = widget.knownMedia?.uid ?? widget.download!.media.knownMediaId
      ..libraryId = widget.download!.media.id;

    httpx
        .withRetry(
          () => widget.apicommunitypublish(
            widget.community?.id ?? uuidx.min(),
            _request,
            options: [authn.request(authn.AuthzCache.meta(context))],
          ),
        )
        .then((_) => widget.onPublished())
        .catchError((cause) {
          setState(() {
            _cause = ds.Errors.httpauto(cause, onTap: _reseterr);
          });
        }, test: httpx.ErrorsTest.httpauto)
        .catchError((cause) {
          setState(() {
            _cause = ds.Error.unknown(cause, onTap: _reseterr);
          });
        })
        .whenComplete(() {
          setState(() {
            _loading = false;
          });
        });
  }

  @override
  Widget build(BuildContext context) {
    final theme = Theme.of(context);
    final defaults = ds.Defaults.of(context);
    final isVideo = mimex.isVideo(widget.download?.media.mimetype ?? '');

    return ds.Loading(
      loading: _loading,
      cause: _cause,
      SingleChildScrollView(
        child: Column(
          mainAxisSize: MainAxisSize.min,
          crossAxisAlignment: CrossAxisAlignment.start,
          children: [
            Text('Ready to Publish', style: theme.textTheme.titleLarge),
            SizedBox(height: defaults.spacing * 2),
            ConfirmationRow(
              label: 'Content',
              value: widget.download?.media.description ?? 'Unknown',
            ),
            ConfirmationRow(
              label: 'Community',
              value: widget.community?.domain ?? 'Unknown',
            ),
            Visibility(
              visible: widget.knownMedia != null,
              child: ConfirmationRow(
                label: 'Media Info',
                value: widget.knownMedia?.description ?? '',
              ),
            ),
            SizedBox(height: defaults.spacing * 2),
            PublishModeEdit(
              publishMode: _request.publishMode,
              onChanged:
                  (mode) => setState(() {
                    _request.publishMode = mode;
                  }),
            ),
            Visibility(
              visible: isVideo,
              child: forms.Checkbox(
                Text('Cross-post to YouTube'),
                value: _request.publishedContent.oauthGoogleId.isNotEmpty,
                onChanged:
                    _oauthGoogleId.isNotEmpty
                        ? (v) => setState(() {
                          _request.publishedContent.oauthGoogleId = (v ?? false) ? _oauthGoogleId : '';
                        })
                        : null,
                description: Text('Upload this content to your linked YouTube account'),
              ),
            ),
            SizedBox(height: defaults.spacing * 2),
            Center(
              child: ElevatedButton(
                onPressed: _loading ? null : _publish,
                child:
                    _loading
                        ? SizedBox(
                          width: 20,
                          height: 20,
                          child: CircularProgressIndicator(strokeWidth: 2),
                        )
                        : Text('Publish'),
              ),
            ),
          ],
        ),
      ),
    );
  }
}

class ConfirmationRow extends StatelessWidget {
  final String label;
  final String value;

  const ConfirmationRow({super.key, required this.label, required this.value});

  @override
  Widget build(BuildContext context) {
    final theme = Theme.of(context);
    final defaults = ds.Defaults.of(context);

    return Padding(
      padding: EdgeInsets.symmetric(vertical: defaults.spacing / 2),
      child: Row(
        crossAxisAlignment: CrossAxisAlignment.start,
        children: [
          SizedBox(
            width: 100,
            child: Text(
              '$label:',
              style: theme.textTheme.bodyMedium?.copyWith(
                fontWeight: FontWeight.bold,
              ),
            ),
          ),
          Expanded(child: Text(value, style: theme.textTheme.bodyMedium)),
        ],
      ),
    );
  }
}
