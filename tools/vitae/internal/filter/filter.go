// Package filter applies a variant: it keeps the sections, roles, highlights and courses
// the variant asks for. It is the only code that understands variants; everything
// downstream (the renderers, the site) reads its output.
package filter

import (
	"fmt"
	"slices"
	"sort"

	"github.com/jonreyna/jonreyna.dev/tools/vitae/internal/model"
)

// DefaultBoost is what each boosted tag adds to a highlight's weight.
const DefaultBoost = 20

// Resolved is a résumé with one variant applied. Sections the variant leaves out are
// empty.
type Resolved struct {
	Variant    string
	Sections   []string
	Basics     model.Basics
	Summary    string
	Experience []model.Role
	Earlier    string
	Skills     []model.SkillGroup
	Education  []model.School
}

// Has reports whether the variant shows section s.
func (r *Resolved) Has(s string) bool { return slices.Contains(r.Sections, s) }

// Apply returns r with the named variant applied. It doesn't validate r; run lint first.
func Apply(r *model.Resume, name string) (*Resolved, error) {
	v, ok := r.Variants[name]
	if !ok {
		return nil, fmt.Errorf("no variant named %q", name)
	}
	out := &Resolved{Variant: name, Sections: v.Sections, Basics: r.Basics}
	if v.Label != "" {
		out.Basics.Label = v.Label
	}
	if out.Has("summary") {
		key := v.Summary
		if key == "" {
			key = "default"
		}
		out.Summary = r.Summary[key]
	}
	if out.Has("experience") {
		out.Experience = roles(r.Experience, v)
	}
	if out.Has("earlier") {
		out.Earlier = r.Earlier
	}
	if out.Has("skills") {
		out.Skills = r.Skills
	}
	if out.Has("education") {
		for _, s := range r.Education {
			s.Courses = courses(s.Courses, v.ShowCourses)
			out.Education = append(out.Education, s)
		}
	}
	return out, nil
}

func roles(all []model.Role, v model.Variant) []model.Role {
	var out []model.Role
	for _, role := range all {
		if v.Since != nil && role.End != nil && role.End.Before(*v.Since) {
			continue
		}
		if v.MaxRoles > 0 && len(out) == v.MaxRoles {
			break
		}
		role.Highlights = highlights(role.Highlights, v)
		out = append(out, role)
	}
	return out
}

// highlights scores each highlight (weight plus boosts), drops excluded and weak ones,
// keeps the best MaxHighlights, and returns the survivors in their original order.
func highlights(all []model.Highlight, v model.Variant) []model.Highlight {
	boostBy := v.BoostBy
	if boostBy == 0 {
		boostBy = DefaultBoost
	}
	type scored struct {
		index, score int
	}
	var kept []scored
	for i, h := range all {
		if anyTag(h.Tags, v.Exclude) {
			continue
		}
		score := h.Weight
		for _, t := range h.Tags {
			if slices.Contains(v.Boost, t) {
				score += boostBy
			}
		}
		if score < v.MinWeight {
			continue
		}
		kept = append(kept, scored{i, score})
	}
	if v.MaxHighlights > 0 && len(kept) > v.MaxHighlights {
		// Highest score first; ties go to the earlier highlight.
		sort.SliceStable(kept, func(a, b int) bool { return kept[a].score > kept[b].score })
		kept = kept[:v.MaxHighlights]
		sort.Slice(kept, func(a, b int) bool { return kept[a].index < kept[b].index })
	}
	out := make([]model.Highlight, 0, len(kept))
	for _, k := range kept {
		out = append(out, all[k.index])
	}
	return out
}

func courses(all []model.Course, show []string) []model.Course {
	if len(show) == 0 {
		return all
	}
	var out []model.Course
	for _, c := range all {
		if anyTag(c.Tags, show) {
			out = append(out, c)
		}
	}
	return out
}

func anyTag(tags, set []string) bool {
	for _, t := range tags {
		if slices.Contains(set, t) {
			return true
		}
	}
	return false
}
