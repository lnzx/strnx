<route lang="yaml">
path: "/settings"
</route>

<template>
  <main class="column bg-color px-5 py-5" id="main">
    <div class="sub-title">Settings</div>

    <div class="box">
      <div class="field is-horizontal">
        <div class="field-label is-normal">
          <label class="label">Account password</label>
        </div>
        <div class="field-body">
          <div class="field">
            <p class="control">
              <input class="input" type="password" placeholder="New password" />
            </p>
          </div>
          <div class="field">
            <p class="control">
              <input
                class="input"
                type="password"
                placeholder="Confirm Password"
              />
            </p>
          </div>
        </div>
      </div>

      <div class="field is-horizontal">
        <div class="field-label"></div>
        <div class="field-body">
          <div class="field is-expanded">
            <div class="field has-addons">
              <p class="control">
                <a class="button is-static"> +44 </a>
              </p>
              <p class="control is-expanded">
                <input
                  class="input"
                  type="tel"
                  placeholder="Your phone number"
                />
              </p>
            </div>
            <p class="help">Do not enter the first zero</p>
          </div>
        </div>
      </div>

      <div class="field is-horizontal">
        <div class="field-label is-normal">
          <label class="label">Department</label>
        </div>
        <div class="field-body">
          <div class="field is-narrow">
            <div class="control">
              <div class="select is-fullwidth">
                <select>
                  <option>Business development</option>
                  <option>Marketing</option>
                  <option>Sales</option>
                </select>
              </div>
            </div>
          </div>
        </div>
      </div>

      <div class="field is-horizontal">
        <div class="field-label">
          <label class="label">Already a member?</label>
        </div>
        <div class="field-body">
          <div class="field is-narrow">
            <div class="control">
              <label class="radio">
                <input type="radio" name="member" />
                Yes
              </label>
              <label class="radio">
                <input type="radio" name="member" />
                No
              </label>
            </div>
          </div>
        </div>
      </div>

      <div class="field is-horizontal">
        <div class="field-label is-normal">
          <label class="label">Ssh password</label>
        </div>
        <div class="field-body">
          <div class="field">
            <div class="control">
              <input class="input" type="password" placeholder="ssh password" />
            </div>
          </div>
        </div>
      </div>

      <div class="field is-horizontal">
        <div class="field-label is-normal">
          <label class="label">Testmail apikey</label>
        </div>
        <div class="field-body">
          <div class="field">
            <div class="control">
              <textarea
                class="textarea"
                placeholder="Enter your testmail apikey"
              >
92d478eca8bc</textarea
              >
            </div>
          </div>
        </div>
      </div>

      <div class="field is-horizontal">
        <div class="field-label">
          <!-- Left empty for spacing -->
        </div>
        <div class="field-body">
          <div class="field">
            <div class="control">
              <button class="button is-link is-fullwidth">Save</button>
            </div>
          </div>
        </div>
      </div>
    </div>
  </main>
</template>
