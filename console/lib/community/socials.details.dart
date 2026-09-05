import 'package:flutter/material.dart';
import 'package:retrovibed/designkit.dart' as ds;
import 'package:retrovibed/httpx.dart' as httpx;
import 'package:retrovibed/authn.dart' as authn;
import 'package:retrovibed/ddisc/plugin/environment.editor.dart';
import 'api.dart';

// Fetches its own catalog + enabled-publisher data for a single community —
// the socials search endpoint backs this expanded details view only, not
// the SocialHome grid itself.
class SocialCommunityDetails extends StatefulWidget {
  final String communityId;
  final FnSocialsSearch search;
  final FnSocialsEnable enable;
  final FnSocialsDisable disable;

  const SocialCommunityDetails({
    super.key,
    required this.communityId,
    required this.search,
    required this.enable,
    required this.disable,
  });

  @override
  State<SocialCommunityDetails> createState() => _SocialCommunityDetailsState();
}

class _SocialCommunityDetailsState extends State<SocialCommunityDetails> with ds.LoadingState {
  SocialsSearchResponse _resp = SocialsSearchResponse();

  void _refresh() {
    setState(() => loading = true);
    httpx
        .withRetry(
          () => widget.search(
            SocialsSearchRequest(limit: ds.Int64(100)),
            options: [authn.request(authn.AuthzCache.meta(context))],
          ),
        )
        .then((response) {
          setState(() {
            _resp = response;
            cause = ds.Error.zero;
          });
        })
        .catchError((cause) {
          setState(() => this.cause = ds.Errors.httpauto(cause, onTap: reseterr));
        }, test: httpx.ErrorsTest.httpauto)
        .catchError((cause) {
          setState(() => this.cause = ds.Error.unknown(cause, onTap: reseterr));
        })
        .whenComplete(() => setState(() => loading = false));
  }

  @override
  void initState() {
    super.initState();
    ds.postframe(_refresh);
  }

  // _configure opens the publisher's .env in the same editor the search
  // plugin list uses; the endpoint hands back the plugin's declared
  // variables with the configured values filled in, so the editor renders
  // a populated form rather than an empty text box.
  Future<void> _configure(BuildContext context, PluginPublisher p) async {
    final modal = ds.modals.of(context);
    if (modal == null) return;

    final auth = [authn.request(authn.AuthzCache.meta(context))];

    modal.push(
      ds.Confirmation.info(
        content: EnvironmentEditor.future(
          p.id,
          publisherenvironment.get(p.id, options: auth),
          onChange: (content) {
            httpx.withRetry(() => publisherenvironment.update(p.id, content, options: auth)).catchError((cause) {
              print("failed to update publisher environment ${cause}");
              return content;
            });
          },
        ),
        done: (_) => modal.push(null),
      ),
    );
  }

  @override
  Widget build(BuildContext context) {
    final social = _resp.items.where((v) => v.community.id == widget.communityId).firstOrNull;
    final enabledIds = (social?.enabled ?? const []).map((e) => e.publisherId).toSet();

    return ds.Loading(
      loading: loading,
      cause: cause,
      Wrap(
        spacing: 8,
        runSpacing: 4,
        children: _resp.catalog.map((p) {
          final enabled = enabledIds.contains(p.id);
          return Row(
            mainAxisSize: MainAxisSize.min,
            children: [
              FilterChip(
                label: Text(p.description.isNotEmpty ? p.description : p.mimetype),
                selected: enabled,
                onSelected: (v) {
                  final auth = [authn.request(authn.AuthzCache.meta(context))];
                  final fut = v
                      ? widget.enable(widget.communityId, p.id, options: auth)
                      : widget.disable(widget.communityId, p.id, options: auth);
                  httpx.withRetry(() => fut).then((_) => _refresh());
                },
              ),
              // the form behind this button is generated from the plugin's
              // own declaration of the variables it understands, so a
              // publisher nobody wrote a settings screen for still gets one.
              ds.LoadingIconButton.edit(
                iconSize: 18.0,
                help: ds.Hint(const Text("configure this publisher")),
                onPressed: () => _configure(context, p),
              ),
            ],
          );
        }).toList(),
      ),
    );
  }
}
