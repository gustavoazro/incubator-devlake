/*
 * Licensed to the Apache Software Foundation (ASF) under one or more
 * contributor license agreements.  See the NOTICE file distributed with
 * this work for additional information regarding copyright ownership.
 * The ASF licenses this file to You under the Apache License, Version 2.0
 * (the "License"); you may not use this file except in compliance with
 * the License.  You may obtain a copy of the License at
 *
 *     http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 *
 */

import styled from '@emotion/styled'

export const Container = styled.div`
  margin-top: 24px;
  padding: 24px;
  font-size: 12px;
  color: #292b3f;
  background-color: #fff;
  box-shadow: 0px 2.4px 4.8px -0.8px rgba(0, 0, 0, 0.1),
    0px 1.6px 8px rgba(0, 0, 0, 0.07);
  border-radius: 8px;

  & > .bp4-button-group {
    display: block;
    margin-top: 16px;
    text-align: right;

    .bp4-button + .bp4-button {
      margin-left: 4px;
    }
  }
`
